package set

type ReadableSet interface {
	Get(key string) interface{}
}

type WritableSet interface {
	Set(key string, value interface{})
}

type CheckableSet interface {
	Has(key string) bool
}

type CheckWritableSet interface {
	WritableSet
	CheckableSet
}

type Set interface {
	ReadableSet
	WritableSet
	CheckableSet
}
