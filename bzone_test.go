package scExplorer

import "github.com/farces/dumb/bufbig"
import (
	"math/big"
	"testing"
)

func TestTotalNumberOfPoints(t *testing.T) {
	checkNumPoints := func(pointsPerSide, dimension int64, expected *big.Int) bool {
		points := bzPoints(pointsPerSide, dimension)
		i := bufbig.NewBigAccumulator()
		for {
			p := <-points
			if p != nil {
				i.AddInt(1)
			} else {
				break
			}
		}
		// x.Cmp(y) returns sgn(x - y)
		return i.Value().Cmp(expected) == 0
	}
	N := int64(64)
	if !checkNumPoints(N, 2, big.NewInt(N*N)) {
		t.Fatalf("Incorrect number of points from bzPoints(%d,%d)", N, 2)
	}
}
