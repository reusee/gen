## Install

```
go get -u github.com/brother-big/gen
```

## Example

```
gen map IntString int string
```

```
type IntStringMap map[int]string

func NewIntStringMap() IntStringMap {
	return make(map[int]string)
}

func (m IntStringMap) Set(key int, value string) {
	m[key] = value
}

func (m IntStringMap) Get(key int) (value string, ok bool) {
	value, ok = m[key]
	return
}

```

```
gen map StrBytes string []byte
```

```
type StrBytesMap map[string][]byte

func NewStrBytesMap() StrBytesMap {
	return make(map[string][]byte)
}

func (m StrBytesMap) Set(key string, value []byte) {
	m[key] = value
}

func (m StrBytesMap) Get(key string) (value []byte, ok bool) {
	value, ok = m[key]
	return
}

```
