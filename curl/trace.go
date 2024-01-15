package curl

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http/httptrace"
	"time"
)

const (
	traceFmt = `TotalTime         : %v
DNSLookupTime     : %v
TCPConnectTime    : %v
TLSHandshakeTime  : %v
FirstResponseTime : %v
ResponseTime      : %v
IsConnReused:     : false
LocalAddr         : %s
RemoteAddr        : %s`
	traceReusedFmt = `TotalTime         : %v
FirstResponseTime : %v
ResponseTime      : %v
IsConnReused:     : true
LocalAddr         : %s
RemoteAddr        : %s`
)

// TraceInfo represents the trace information.
type TraceInfo struct {
	// DNSLookupTime is a duration that transport took to perform
	// DNS lookup.
	DNSLookupTime time.Duration

	// ConnectTime is a duration that took to obtain a successful connection.
	ConnectTime time.Duration

	// TCPConnectTime is a duration that took to obtain the TCP connection.
	TCPConnectTime time.Duration

	// TLSHandshakeTime is a duration that TLS handshake took place.
	TLSHandshakeTime time.Duration

	// FirstResponseTime is a duration that server took to respond first byte since
	// connection ready (after tls handshake if it's tls and not a reused connection).
	FirstResponseTime time.Duration

	// ResponseTime is a duration since first response byte from server to
	// request completion.
	ResponseTime time.Duration

	// TotalTime is a duration that total request took end-to-end.
	TotalTime time.Duration

	// IsConnReused is whether this connection has been previously
	// used for another HTTP request.
	IsConnReused bool

	// IsConnWasIdle is whether this connection was obtained from an
	// idle pool.
	IsConnWasIdle bool

	// ConnIdleTime is a duration how long the connection was previously
	// idle, if IsConnWasIdle is true.
	ConnIdleTime time.Duration

	// LocalAddr returns the local network address.
	LocalAddr net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr net.Addr
}

// String return the details of trace information.
func (t TraceInfo) String() string {
	var localAddr, remoteAddr string
	if t.LocalAddr != nil {
		localAddr = t.LocalAddr.String()
	}
	if t.RemoteAddr != nil {
		remoteAddr = t.RemoteAddr.String()
	}
	if t.IsConnReused {
		return fmt.Sprintf(traceReusedFmt, t.TotalTime, t.FirstResponseTime, t.ResponseTime, localAddr, remoteAddr)
	}
	return fmt.Sprintf(traceFmt, t.TotalTime, t.DNSLookupTime, t.TCPConnectTime, t.TLSHandshakeTime,
		t.FirstResponseTime, t.ResponseTime, localAddr, remoteAddr)
}

type clientTrace struct {
	getConn              time.Time
	dnsStart             time.Time
	dnsDone              time.Time
	connectDone          time.Time
	tlsHandshakeStart    time.Time
	tlsHandshakeDone     time.Time
	gotConn              time.Time
	gotFirstResponseByte time.Time
	endTime              time.Time
	gotConnInfo          httptrace.GotConnInfo
}

func (t *clientTrace) createContext(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(
		ctx,
		&httptrace.ClientTrace{
			DNSStart: func(_ httptrace.DNSStartInfo) {
				t.dnsStart = time.Now()
			},
			DNSDone: func(_ httptrace.DNSDoneInfo) {
				t.dnsDone = time.Now()
			},
			ConnectStart: func(_, _ string) {
				if t.dnsDone.IsZero() {
					t.dnsDone = time.Now()
				}
				if t.dnsStart.IsZero() {
					t.dnsStart = t.dnsDone
				}
			},
			ConnectDone: func(net, addr string, err error) {
				t.connectDone = time.Now()
			},
			GetConn: func(_ string) {
				t.getConn = time.Now()
			},
			GotConn: func(ci httptrace.GotConnInfo) {
				t.gotConn = time.Now()
				t.gotConnInfo = ci
			},
			GotFirstResponseByte: func() {
				t.gotFirstResponseByte = time.Now()
			},
			TLSHandshakeStart: func() {
				t.tlsHandshakeStart = time.Now()
			},
			TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
				t.tlsHandshakeDone = time.Now()
			},
		},
	)
}
