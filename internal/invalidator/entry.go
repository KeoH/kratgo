package invalidator

import "sync"

var entryPool = sync.Pool{
	New: func() interface{} {
		return new(Entry)
	},
}

// AcquireEntry ...
func AcquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

// ReleaseEntry ...
func ReleaseEntry(e *Entry) {
	e.Reset()
	entryPool.Put(e)
}

// Reset ...
func (h *Header) Reset() {
	h.Key = ""
	h.Value = ""
}

// Reset ...
func (e *Entry) Reset() {
	e.Action = ""
	e.Host = ""
	e.Path = ""

	e.Header.Reset()
}
