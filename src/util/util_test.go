package util

import "testing"

func TestDown(t *testing.T) {
	DownFile("https://czqu.github.io/", "./test")
}
