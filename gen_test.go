package main

import (
	"testing"
)

func TestMap(t *testing.T) {
	src, err := Gen("map", "IntString", "int", "string")
	if err != nil {
		t.Fatal(err)
	}
	if src != `type IntStringMap map[int]string
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
` {
		t.Fatal()
	}
}
