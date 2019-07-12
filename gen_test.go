package main

import (
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	src, err := Gen("map", "IntString", "int", "string")
	if err != nil {
		t.Fatal(err)
	}
	println(strings.Repeat("-", 32))
	println(src)
}
