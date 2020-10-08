package lego

// deprecated
type CloserGroup struct{}

// deprecated
func (CloserGroup) Close() error { return nil }

// deprecated
func (CloserGroup) Add(interface{}) {}
