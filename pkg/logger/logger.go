package logger

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// InitLogger initializes the logger based on the environment
func InitLogger(isDebug bool) {
	var err error
	if isDebug {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	sugar = logger.Sugar()
}

// GetLogger returns the initialized logger
func GetLogger() *zap.Logger {
	return logger
}

// GetSugar returns the initialized sugared logger
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// LogRequest logs HTTP request details and updates Prometheus metrics
func LogRequest(r *http.Request) func() {
	start := time.Now()
	sugar.Infow("Received request",
		"method", r.Method,
		"path", r.URL.Path,
		"remoteAddr", r.RemoteAddr,
	)

	httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

	return func() {
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	}
}

// MetricsHandler returns the Prometheus metrics handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
