package synse

import (
	"github.com/davidkleiven/gosfft/sfft"
	"gonum.org/v1/gonum/mat"
)

// do_convolve, get_psf_freq, get_psf_spat,

func getGrad(M, N int) *mat.CDense {
	gradKernel := mat.NewDense(M, N, nil)
	gradKernel.Set(M/2, N/2+1, 1)
	gradKernel.Set(M/2, N/2, -1)
	arr := spatToFreq(gradKernel)
	// sfft.Center2(arr)
	return arr
}

func getConvolver(M, N int, a1, a2, l1, l2, s float64) *mat.CDense {
	gradMask := getGrad(M, N)
	psfMask := GetPsfFreq(M, N, a1, a2, l1, l2, s)
	convolved := mat.NewCDense(M, N, nil)
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			convolved.Set(i, j, psfMask.At(i, j)*gradMask.At(i, j))
		}
	}
	return convolved
}

func GetPsfFreq(M, N int, a1, a2, l1, l2, s float64) *mat.CDense {
	psf := GetPSF(M, N, a1, a2, l1, l2, s)
	complexPsf := complexify(psf)
	sfft.Center2(complexPsf)
	return complexPsf
}

func GetPsfSpat(M, N int, a1, a2, l1, l2, s float64) *mat.CDense {
	psfFreq := GetPsfFreq(M, N, a1, a2, l1, l2, s)
	psfSpat := freqToSpat(psfFreq)
	return psfSpat
}

func DoConvolve(data *mat.Dense, a1, a2, l1, l2, s float64) *mat.CDense {
	M, N := data.Dims()
	convolver := getConvolver(M, N, a1, a2, l1, l2, s)
	F := spatToFreq(data)
	maskedF := mat.NewCDense(M, N, nil)
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			maskedF.Set(i, j, F.At(i, j)*convolver.At(i, j))
		}
	}
	return freqToSpat(maskedF)
}

func freqToSpat(data *mat.CDense) *mat.CDense {
	freq := ifft2(data)
	return freq
}

func spatToFreq(data *mat.Dense) *mat.CDense {
	complexData := complexify(data)
	img := fftn(complexData)
	return img
}
