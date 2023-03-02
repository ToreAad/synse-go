package synse

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"

	"gonum.org/v1/gonum/interp"
	"gonum.org/v1/gonum/mat"
)

type SeismicColor struct {
	redInt   *interp.PiecewiseLinear
	blueInt  *interp.PiecewiseLinear
	greenInt *interp.PiecewiseLinear
}

func NewSeismicColor() *SeismicColor {
	xs := []float64{0, 0.25, 0.5, 0.75, 1.0}
	red := []float64{0, 0, 1, 0.8314, 0.5}
	green := []float64{0, 0.375, 1, 0.375, 0}
	blue := []float64{0.5, 0.8314, 1, 0, 0}

	redInt := &interp.PiecewiseLinear{}
	redInt.Fit(xs, red)
	greenInt := &interp.PiecewiseLinear{}
	greenInt.Fit(xs, green)
	blueInt := &interp.PiecewiseLinear{}
	blueInt.Fit(xs, blue)
	return &SeismicColor{
		redInt:   redInt,
		blueInt:  blueInt,
		greenInt: greenInt,
	}
}

func (sc *SeismicColor) Color(z float64) color.RGBA {
	red := uint8(sc.redInt.Predict(z) * 255)
	green := uint8(sc.greenInt.Predict(z) * 255)
	blue := uint8(sc.blueInt.Predict(z) * 255)
	return color.RGBA{R: red, G: green, B: blue, A: 0xFF}
}

func toSeismicImage(arr *mat.Dense) image.Image {
	m, n := arr.Dims()
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{m, n}})
	mx := mat.Max(arr)
	mn := mat.Min(arr)

	if mx == mn {
		return img
	}

	sc := NewSeismicColor()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			z := (arr.At(i, j) - mn) / (mx - mn)
			img.Set(i, j, sc.Color(z))
		}
	}
	return img
}

func toGrayImage(arr *mat.Dense) image.Image {
	m, n := arr.Dims()
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{m, n}})
	mx := mat.Max(arr)
	mn := mat.Min(arr)
	if mx == mn {
		return img
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			z := uint8(255 * (arr.At(i, j) - mn) / (mx - mn))
			c := color.Gray{Y: z}
			img.Set(i, j, c)
		}
	}
	return img
}

func realValue(c color.Color) (float64, error) {
	val, ok := color.GrayModel.Convert(c).(color.Gray)
	if !ok {
		return 0.0, errors.New("Failed to convert image point to grayscale")
	}
	return float64(val.Y) / 255.0, nil
}

func toMat(im image.Image) (*mat.Dense, error) {
	bound := im.Bounds()
	m := bound.Max.X - bound.Min.X
	n := bound.Max.Y - bound.Min.Y
	arr := mat.NewDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			col := im.At(i+bound.Min.X, j+bound.Min.Y)
			val, err := realValue(col)
			if err != nil {
				return nil, err
			}
			arr.Set(i, j, val)
		}
	}
	return arr, nil
}

func fromDataUrl(dataUrl string) (image.Image, error) {
	data := strings.Split(dataUrl, "data:image/png;base64,")
	if len(data) < 2 {
		return nil, errors.New("could not parse dataurl")
	}
	binData, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.Write(binData)
	img, err := png.Decode(buf)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func toDataUrl(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", data), nil
}
