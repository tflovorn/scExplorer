package solve

// Three-point algorithm for second derivative of f(x) at x0. Error ~ h^2.
func Simple2ndDiff(f Func1D, x0, h float64) (float64, error) {
	fp, err := f(x0 + h)
	if err != nil {
		return 0.0, err
	}
	f0, err := f(x0)
	if err != nil {
		return 0.0, err
	}
	fm, err := f(x0 - h)
	if err != nil {
		return 0.0, err
	}
	return (fp - 2.0*f0 + fm) / (h * h), nil
}

// Four-point algorithm for mixed partial second derivative. Error ~ hf^2*hg^2.
func SimpleMixed2ndDiff(f func(x, y float64) (float64, error), x0, y0, hx, hy float64) (float64, error) {
	fpp, err := f(x0+hx, y0+hy)
	if err != nil {
		return 0.0, err
	}
	fpm, err := f(x0+hx, y0-hy)
	if err != nil {
		return 0.0, err
	}
	fmp, err := f(x0-hx, y0+hy)
	if err != nil {
		return 0.0, err
	}
	fmm, err := f(x0-hx, y0-hy)
	if err != nil {
		return 0.0, err
	}
	return (fpp - fpm - fmp + fmm) / (4.0 * hx * hy), nil
}
