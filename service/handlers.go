package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xnyo/lineameteo-prometheus/lineameteo"
	"net/http"
)

const weatherAPIURL = "https://retemeteo.lineameteo.it/rete_json.json"

func getMetrics(w http.ResponseWriter, r *http.Request) {
	ic := GetInnerContext(r.Context())

	// TODO: move to lineameteo.go
	// Prepare request
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, weatherAPIURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("weather request: %s", err), http.StatusServiceUnavailable)
		return
	}

	// 403 otherwise
	req.Header.Set("Referer", "https://retemeteo.lineameteo.it/index.php")

	// Do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("weather api error: %s", err), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// The output is not JSON, what the hell
	// Output of type:
	// markers = JSON
	// Keep reading until "="
	reader := bufio.NewReader(resp.Body)
	for {
		b, _, err := reader.ReadRune()
		if err != nil {
			http.Error(w, fmt.Sprintf("read rune: %s", err), http.StatusServiceUnavailable)
			return
		}
		if b == '=' {
			break
		}
	}

	// Pass the rest to the JSON parser
	var locations []lineameteo.Location
	if err := json.NewDecoder(reader).Decode(&locations); err != nil {
		http.Error(w, fmt.Sprintf("json decode: %s", err), http.StatusInternalServerError)
		return
	}

	// Find wanted locations
	for _, loc := range locations {
		if gauge, ok := ic.Service.PrometheusGauges[loc.ID]; ok {
			// log.Println(loc.ID, loc.Name, loc.Temperatura)

			// Set all gauges
			gauge.Temperature.Set(float64(loc.Temperatura))
			gauge.MaxTemperature.Set(float64(loc.Minima))
			gauge.MinTemperature.Set(float64(loc.Massima))
			gauge.Humidity.Set(float64(loc.Umidita))
			gauge.Pressure.Set(float64(loc.Pressione))
		}
	}

	// Call prometheus handler
	ic.Service.prometheusHandler.ServeHTTP(w, r)
}

func init() {
	router.Group(func(r chi.Router) {
		r.Use(middleware.Throttle(4))
		r.Get("/metrics", getMetrics)
	})
}
