package lineameteo

import (
	"strconv"
	"strings"
)

type Location struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Lat            string  `json:"lat"`
	Lng            string  `json:"lng"`
	Alt            string  `json:"alt"`
	Webcam         string  `json:"webcam"`
	Webcam2        string  `json:"webcam2"`
	Time           string  `json:"time"`
	Temperatura    Float64 `json:"temperatura"`
	Umidita        Float64 `json:"umidita"`
	Pressione      Float64 `json:"pressione"`
	Vento          string  `json:"vento"`
	VentoDirezione string  `json:"vento_direzione"`
	Minima         Float64 `json:"minima"`
	Massima        Float64 `json:"massima"`
	Raffica        string  `json:"raffica"`
	Pioggia        string  `json:"pioggia"`
	RainRate       string  `json:"rain_rate"`
	Modello        string  `json:"modello"`
	Url            string  `json:"url"`
	Condizioni     string  `json:"condizioni"`
}

type Float64 float64

func (f *Float64) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	nt, err := strconv.ParseFloat(s, 64)
	if err != nil {
		*f = 0.0
		return nil
	}
	*f = Float64(nt)
	return nil
}
