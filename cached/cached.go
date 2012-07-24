// Provide a cached interface to stateless math functions (sin, cos, etc.)
package cached

import "math"

var cosCache func(float64) float64
var sinCache func(float64) float64

// Return a wrapper over fn which returns a cached value when available.
func WrapMath(fn func(float64) float64) func(float64) float64 {
	cache := make(map[float64]float64)
	return func(x float64) float64 {
		val, ok := cache[x]
		if !ok {
			val = fn(x)
			cache[x] = val
		}
		return val
	}
}

// Cached version of math.Cos.
func Cos(x float64) float64 {
	// could remove this check by forcing in Init() call before this
	if cosCache == nil {
		cosCache = WrapMath(math.Cos)
	}
	return cosCache(x)
}

// Cached version of math.Sin.
func Sin(x float64) float64 {
	if sinCache == nil {
		sinCache = WrapMath(math.Sin)
	}
	return sinCache(x)
}

// Possible functin to wrap bzPoints:
// Return a caching wrapper for the channel-producing function fn.
// fn should return a finite number of values, then close the channel.
// func WrapChan(fn func(int, int) <-chan vec.Vector)
