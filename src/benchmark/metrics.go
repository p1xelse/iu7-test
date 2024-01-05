package benchmark

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	encodeAlgorithmLbl = "encoding_algorithm"
	decodeAlgorithmLbl = "decoding_algorithm"
)

//go:generate mockgen --source=metrics.go --destination=mocks/metrics.go -package=mocks
type Metrics interface {
	EncodeRecordTime(encodingAlgorithm string, duration float64)
	DecodeRecordTime(decodingAlgorithm string, duration float64)
}

type benchMetrics struct {
	encodeTime *prometheus.HistogramVec
	decodeTime *prometheus.HistogramVec
}

func (bm *benchMetrics) EncodeRecordTime(encodingAlgorithm string, duration float64) {
	bm.encodeTime.WithLabelValues(encodingAlgorithm).Observe(duration)
}

func (bm *benchMetrics) DecodeRecordTime(decodingAlgorithm string, duration float64) {
	bm.decodeTime.WithLabelValues(decodingAlgorithm).Observe(duration)
}

func NewMetrics() *benchMetrics {
	encodeTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "encode_duration_nanoseconds",
		Help: "Time taken for encoding in nanoseconds",
		Buckets: prometheus.LinearBuckets(500, 200, 20),
	}, []string{encodeAlgorithmLbl})

	// Инициализация метрик для времени декода
	decodeTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "decode_duration_nanoseconds",
		Help: "Time taken for decoding in nanoseconds",
		Buckets: prometheus.LinearBuckets(500, 200, 20),
	}, []string{decodeAlgorithmLbl})

	return &benchMetrics{
		encodeTime: encodeTime,
		decodeTime: decodeTime,
	}
}
