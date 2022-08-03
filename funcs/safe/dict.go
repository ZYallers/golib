package safe

import "sync"

type Dict struct {
	data map[string]interface{}
	*sync.RWMutex
}

func NewDict() *Dict {
	return &Dict{map[string]interface{}{}, &sync.RWMutex{}}
}

func (d *Dict) Len() int {
	d.RLock()
	defer d.RUnlock()
	return len(d.data)
}

func (d *Dict) Data() map[string]interface{} {
	d.RLock()
	defer d.RUnlock()
	cp := make(map[string]interface{}, len(d.data))
	for k, v := range d.data {
		cp[k] = v
	}
	return cp
}

func (d *Dict) Get(key string) (interface{}, bool) {
	d.RLock()
	defer d.RUnlock()
	oldValue, ok := d.data[key]
	return oldValue, ok
}

func (d *Dict) GetOrPut(key string, value interface{}) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	if !ok {
		d.data[key] = value
	}
	return oldValue, ok
}

func (d *Dict) GetOrPutFunc(key string, f func(string) (interface{}, error)) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	if ok {
		return oldValue, ok
	}
	if v, err := f(key); err == nil {
		d.data[key] = v
		return v, ok
	}
	return nil, ok
}

func (d *Dict) Put(key string, value interface{}) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	d.data[key] = value
	return oldValue, ok
}

func (d *Dict) Delete(key string) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	if ok {
		delete(d.data, key)
	}
	return oldValue, ok
}
