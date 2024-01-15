package logger

type config struct {
	maxSize    int
	maxBackups int
	maxAge     int
	localTime  bool
	compress   bool
}

type Option func(c *config)

func WithMaxSize(s int) Option {
	return func(c *config) {
		c.maxSize = s
	}
}

func WithMaxAge(s int) Option {
	return func(c *config) {
		c.maxAge = s
	}
}

func WithMaxBackups(b int) Option {
	return func(c *config) {
		c.maxBackups = b
	}
}

func WithLocalTime(b bool) Option {
	return func(c *config) {
		c.localTime = b
	}
}

func WithCompress(b bool) Option {
	return func(c *config) {
		c.compress = b
	}
}
