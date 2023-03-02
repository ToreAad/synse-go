package synse

// Reference https://github.com/davidkleiven/gosfft/blob/master/sfft/fft.go

import (
	"github.com/davidkleiven/gosfft/sfft"
	"gonum.org/v1/gonum/mat"
)

func ifftshift(data *mat.CDense) *mat.CDense {
	m, n := data.Dims()
	centeredData := mat.NewCDense(m, n, data.RawCMatrix().Data)
	sfft.Center2(centeredData)
	return centeredData
}

func ifft2(data *mat.CDense) *mat.CDense {
	m, n := data.Dims()
	ft := sfft.NewFFT2(m, n)
	ftData := ft.IFFT(data.RawCMatrix().Data)
	ftMat := mat.NewCDense(m, n, ftData)
	return ftMat
}

func realify(data *mat.CDense) *mat.Dense {
	m, n := data.Dims()
	realData := mat.NewDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			realData.Set(i, j, real(data.At(i, j)))
		}
	}
	return realData
}

func complexify(data *mat.Dense) *mat.CDense {
	m, n := data.Dims()
	cmplxData := mat.NewCDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			cmplxData.Set(i, j, complex(data.At(i, j), 0))
		}
	}
	return cmplxData
}

func fftn(data *mat.CDense) *mat.CDense {
	m, n := data.Dims()
	ft := sfft.NewFFT2(m, n)
	ftData := ft.FFT(data.RawCMatrix().Data)
	ftMat := mat.NewCDense(m, n, ftData)
	return ftMat
}
