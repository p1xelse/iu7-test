package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/mailru/easyjson"
)

type encodeBenchmarker struct {
	entryBuilder *testutils.EntryBuilder
	tagBuilder   *testutils.TagBuilder

	benchMetrics Metrics
}

func NewEncodeBenchmarker(metrics Metrics) *encodeBenchmarker {
	return &encodeBenchmarker{
		benchMetrics: metrics,
		entryBuilder: testutils.NewEntryBuilder(),
		tagBuilder:   testutils.NewTagBuilder(),
	}
}

func (db *encodeBenchmarker) getDummyEntry() models.Entry {
	entry := db.entryBuilder.
		WithID(1).
		WithDescription("entry").
		WithUserID(1).
		WithTagList([]models.Tag{db.tagBuilder.WithName("hello").Build()}).
		Build()

	return entry
}

func (db *encodeBenchmarker) benchEncode_EncodingJson(n int) func(b *testing.B) {
	entry := db.getDummyEntry()
	return func(b *testing.B) {
		for j := 0; j < n; j++ {
			for i := 0; i < b.N; i++ {
				startTime := time.Now()
				_, err := json.Marshal(entry)
				if err != nil {
					panic(err)
				}
				elapsed := time.Since(startTime)
				db.benchMetrics.EncodeRecordTime(encodingJsonLblValue, float64(elapsed.Nanoseconds()))
			}
		}
	}
}

func (db *encodeBenchmarker) benchEncode_EasyJson(n int) func(b *testing.B) {
	entry := db.getDummyEntry()
	return func(b *testing.B) {
		for j := 0; j < n; j++ {
			for i := 0; i < b.N; i++ {
				startTime := time.Now()
				_, err := easyjson.Marshal(entry)
				if err != nil {
					panic(err)
				}
				elapsed := time.Since(startTime)
				db.benchMetrics.EncodeRecordTime(goJsonLblValue, float64(elapsed.Nanoseconds()))
			}
		}
	}
}

func (db *encodeBenchmarker) EncodeBenchmark(n int) (res []string, err error) {
	encodeEncodingJson := db.benchEncode_EncodingJson(n)
	encodeGoJson := db.benchEncode_EasyJson(n)

	resultsEncodingJson := testing.Benchmark(encodeEncodingJson)
	res = append(res, fmt.Sprintf("encoding/json.Marshall -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op",
		resultsEncodingJson.N, resultsEncodingJson.NsPerOp(), resultsEncodingJson.AllocsPerOp(), resultsEncodingJson.AllocedBytesPerOp(),
	))

	resultsGoJson := testing.Benchmark(encodeGoJson)
	res = append(res, fmt.Sprintf("easy-json.Marshall       -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsGoJson.N, resultsGoJson.NsPerOp(), resultsGoJson.AllocsPerOp(), resultsGoJson.AllocedBytesPerOp(),
	))

	return res, nil
}
