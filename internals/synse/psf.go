package synse

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func gaussian(x, mu, sig float64) float64 {
	return math.Exp(-math.Pow(x-mu, 2) / (2 * math.Pow(sig, 2.0)))
}

func getGaussian(x, dev float64) float64 {
	if dev == 0 {
		dev = 1e-4
	}
	mn := gaussian(0, 0.5, dev)
	val := gaussian(x, 0.5, dev)
	mx := 1.0
	result := (val - mn) / (mx - mn)
	return result
}

func rEclipse(a, b, o float64) float64 {
	return (a * b) / math.Pow((math.Pow((b*math.Cos(o)), 2)+math.Pow((a*math.Sin(o)), 2)), 0.5)
}

func transformAngle(l, o float64) float64 {
	y1 := math.Sin(o)
	x1 := math.Cos(o)
	return math.Atan2(l*y1, x1)
}

func GetPSF(M, N int, a1, a2, l1, l2, s float64) *mat.Dense {
	mask := mat.NewDense(M, N, nil)
	a1 -= 90
	a2 -= 90
	for x := -M / 2; x < M/2; x++ {
		for y := -N / 2; y < N/2; y++ {
			i_x := x + M/2
			i_y := y + N/2
			mxl1 := math.Min(l1*float64(M)*0.5, float64(M)*0.5)
			mxl2 := math.Min(l2*mxl1, float64(N)*0.5)
			fx := float64(x)
			fy := float64(y)
			r := math.Sqrt(fx*fx + fy*fy)
			o := -math.Atan2(fx, fy)
			maxLen := rEclipse(mxl1, mxl2, o)
			o1 := transformAngle(l2, a1*math.Pi/180.0)
			o2 := transformAngle(l2, a2*math.Pi/180.0)

			if (r < maxLen) && (o > o1) && (o < o2) {
				dev := 0.25 * (1 - s)
				val := getGaussian(r/maxLen, dev)
				if s == 0 {
					mask.Set(i_x, N-i_y, val)
				} else if math.Abs(val) > 1e-1 {
					mask.Set(i_x, N-i_y, val)
				}
			}
		}
	}
	return mask
}
