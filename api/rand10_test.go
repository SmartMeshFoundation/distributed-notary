package api_test

import (
	"math/rand"
	"testing"
)

func rand7() int {
	return rand.Int()%7 + 1
}
func rand10() int {
	var m [7][7]bool
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if i*7+j <= 39 {
				m[i][j] = true
			}
		}
	}
	for {
		var r = rand7() - 1
		var c = rand7() - 1
		if m[r][c] {
			return (r*7+c)/4 + 1
		}
	}

}

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("t=%d", rand10())
	}
}
