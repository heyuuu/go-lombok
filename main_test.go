package main

import (
	"strings"
	"testing"
)

func TestGenFunc(t *testing.T) {
	args := strings.Split("./gp-gen -cmd generate -d tmp/code", " ")
	run(args)
}
