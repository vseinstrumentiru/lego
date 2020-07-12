package set

func NewSimple() Set {
	return make(simple)
}

type simple map[string]interface{}

func (s simple) Has(key string) bool {
	_, ok := s[key]

	return ok
}

func (s simple) Get(key string) interface{} {
	v, ok := s[key]
	if !ok {
		return nil
	}
	return v
}

func (s simple) Set(key string, value interface{}) {
	s[key] = value
}
