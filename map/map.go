package m

type Key interface{}

type Value interface{}

// name: {Name}Map
type Map map[Key]Value

// name: New{Name}Map
func New() Map {
	return make(map[Key]Value)
}

func (m Map) Set(key Key, value Value) {
	m[key] = value
}

func (m Map) Get(key Key) (value Value, ok bool) {
	value, ok = m[key]
	return
}
