package parallel

import (
	"testing"
)

func TestRunCounting(t *testing.T) {
	N := 8
	state := make([]int, N)
	F := func(i int, resp chan<- error) {
		state[i] = i
		resp <- nil
	}
	errs := Run(F, N)
	for i, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
		if state[i] != i {
			t.Fatal("incorrect state")
		}
	}
}
