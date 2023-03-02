package synse

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/davidkleiven/gosfft/sfft"
)

type synseRequestModel struct {
	Angle1    float64 `json:"angle1"`
	Angle2    float64 `json:"angle2"`
	Length    float64 `json:"length"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Lithology string  `json:"lithology"`
}

type SynseResponseModel struct {
	SeismicImage string `json:"seismic_image"`
}

func SeismicHandler(w http.ResponseWriter, r *http.Request) {
	var synseRequest synseRequestModel
	json.NewDecoder(r.Body).Decode(&synseRequest)
	strataImg, err := fromDataUrl(synseRequest.Lithology)
	a1 := synseRequest.Angle1
	a2 := synseRequest.Angle2
	l1 := synseRequest.Length
	l2 := 1.0
	s := 0.0

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}
	strataData, err := toMat(strataImg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}
	convolvedImage := DoConvolve(strataData, a1, a2, l1, l2, s)
	centeredConvolvedImage := ifftshift(convolvedImage)
	realCenteredConvolvedImage := realify(centeredConvolvedImage)

	img := toSeismicImage(realCenteredConvolvedImage)
	dataUrl, err := toDataUrl(img)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}

	synseResponse := SynseResponseModel{SeismicImage: dataUrl}
	json.NewEncoder(w).Encode(synseResponse)
}

type TemporalResponseModel struct {
	PsfImageFrequency string `json:"psf_image_frequency"`
}

func PsfTemporalHandler(w http.ResponseWriter, r *http.Request) {
	var synseRequest synseRequestModel
	json.NewDecoder(r.Body).Decode(&synseRequest)
	M := synseRequest.Width
	N := synseRequest.Height
	a1 := synseRequest.Angle1
	a2 := synseRequest.Angle2
	l1 := synseRequest.Length
	l2 := 1.0
	s := 0.0

	psfFreq := GetPsfFreq(M, N, a1, a2, l1, l2, s)
	sfft.Center2(psfFreq)
	psfFreqReal := realify(psfFreq)
	img := toGrayImage(psfFreqReal)
	dataUrl, err := toDataUrl(img)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}

	temporalResponse := TemporalResponseModel{PsfImageFrequency: dataUrl}
	json.NewEncoder(w).Encode(temporalResponse)
}

type SpatialresponseModel struct {
	PsfImageSpatial string `json:"psf_image_spatial"`
}

func PsfSpatialHandler(w http.ResponseWriter, r *http.Request) {
	var synseRequest synseRequestModel
	json.NewDecoder(r.Body).Decode(&synseRequest)
	M := synseRequest.Width
	N := synseRequest.Height
	a1 := synseRequest.Angle1
	a2 := synseRequest.Angle2
	l1 := synseRequest.Length
	l2 := 1.0
	s := 0.0

	psfSpat := GetPsfSpat(M, N, a1, a2, l1, l2, s)
	psfSpatCentered := ifftshift(psfSpat)
	realPsfSpat := realify(psfSpatCentered)
	img := toSeismicImage(realPsfSpat)
	dataUrl, err := toDataUrl(img)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}

	spatialResponse := SpatialresponseModel{PsfImageSpatial: dataUrl}
	json.NewEncoder(w).Encode(spatialResponse)
}
