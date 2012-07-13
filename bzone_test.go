package scExplorer

import "github.com/farces/dumb/bufbig"
import (
	"fmt"
	"math/big"
	"testing"
)

func TestTotalNumberOfPoints(t *testing.T) {
	checkNumPoints := func(pointsPerSide, dimension int, expected *big.Int) (*big.Int, bool) {
		points := bzPoints(pointsPerSide, dimension)
		i := bufbig.NewBigAccumulator()
		for {
			p := <-points
			//fmt.Printf("p = %v\n", p)
			if p != nil {
				i.AddInt(1)
			} else {
				break
			}
		}
		iv := i.Value()
		// x.Cmp(y) returns sgn(x - y)
		return iv, iv.Cmp(expected) == 0
	}
	N := 64
	count, correct := checkNumPoints(N, 2, big.NewInt(int64(N*N)))
	if !correct {
		msg := fmt.Sprintf("Incorrect number of points from bzPoints(%d,%d)\n", N, 2)
		msg += fmt.Sprintf("Got %d, expected %d\n", count, int64(N*N))
		t.Fatal(msg)
	}
}
