package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xnyo/lineameteo-prometheus/metrics"
	"net/http"
	"strings"
	"time"
)

const httpRequestTimeout = time.Second * 60

// router is global router that will be mounted on "/"
// This can be used to register handlers globally in init()
var router = chi.NewRouter()

type Service struct {
	HTTPServer *http.Server

	PrometheusGauges    metrics.GaugeMap
	PrometheusCollector prometheus.Gatherer

	prometheusHandler http.Handler

	r *chi.Mux
}

type InnerContext struct {
	Service *Service
}

func GetInnerContext(ctx context.Context) *InnerContext {
	return ctx.Value("innerContext").(*InnerContext)
}

func (s *Service) innerContextMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Inject InnerContext into the Context
			ctx := context.WithValue(r.Context(), "innerContext", &InnerContext{
				Service: s,
			})
			// Call handler
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func NewService(addr string, wantedIDs []string) *Service {
	// Create gauges for each wanted location
	reg := prometheus.NewRegistry()
	gauges := make(metrics.GaugeMap, len(wantedIDs))
	for _, id := range wantedIDs {
		if id == "" {
			continue
		}
		id = strings.TrimSpace(id)
		g := metrics.NewLocationGauges(id)
		g.MustRegister(reg)
		gauges[id] = g
	}

	// Create service
	s := Service{
		r:                   chi.NewRouter(),
		PrometheusGauges:    gauges,
		PrometheusCollector: reg,
		prometheusHandler:   promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	}

	// Middleware stack
	s.r.Use(middleware.RequestID)
	s.r.Use(middleware.RealIP)
	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Recoverer)
	s.r.Use(middleware.Timeout(httpRequestTimeout))
	s.r.Use(s.innerContextMiddleware())
	s.r.Mount("/", router)

	// std HTTP server
	s.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	return &s
}
