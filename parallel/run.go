package parallel

import (
	"fmt"
	"runtime"
	"time"
)

// Execute F(x, cerr) for x = 0..N-1.
func Run(F func(int, chan<- error), N int) []error {
	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)
	numActive, i := 0, 0
	resp := make([]chan error, ncpu)
	fmt.Println("ncpu", ncpu)
	respOwner := make([]int, ncpu)
	for i := 0; i < ncpu; i++ {
		resp[i] = make(chan error)
		respOwner[i] = -1 // -1 == no owner of this channel
	}
	errs := make([]error, N)
	// iterate untill all F's are launched and complete
	for !(numActive == 0 && i >= N) {
		if numActive < ncpu && i < N {
			// spawn a new process
			index := findEmptyChan(respOwner)
			go F(i, resp[index])
			respOwner[index] = i
			numActive += 1
			i += 1
		} else {
			// wait for a reply and handle it
			reply, index := waitOn(resp)
			errs[respOwner[index]] = reply
			respOwner[index] = -1
			numActive -= 1
		}
	}
	return errs
}

// Get the first unused channel in resp
func findEmptyChan(respOwner []int) int {
	for i := 0; i < len(respOwner); i++ {
		if respOwner[i] == -1 {
			return i
		}
	}
	return -1 // shouldn't get here; will result in panic in Run()
}

func waitOn(resp []chan error) (error, int) {
	// wait for the first response
	for {
		for i := 0; i < len(resp); i++ {
			select {
			case err := <-resp[i]:
				return err, i
			default:
				continue
			}
		}
		// sleep to avoid busy waiting
		sleepFor := 50 * time.Millisecond
		time.Sleep(sleepFor)
	}
	return nil, 0.0
}
