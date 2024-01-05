package benchmark

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	integrationutils "timetracker/integration_test/integration_utils"
	"timetracker/models/dto"

	"github.com/mailru/easyjson"
)

type decodeBenchmarker struct {
	entryReqBuilder *integrationutils.ReqCreateUpdateEntryBuilder

	benchMetrics Metrics
}

func NewDecodeBenchmarker(metrics Metrics) *decodeBenchmarker {
	return &decodeBenchmarker{
		benchMetrics:    metrics,
		entryReqBuilder: integrationutils.NewReqCreateUpdateEntryBuilder(),
	}
}

func (db *decodeBenchmarker) getDummyEntryReq() ([]byte, error) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entryDto, err := db.entryReqBuilder.
		WithDescription("entry").
		WithTimeStart(startTime).
		WithTimeEnd(startTime.Add(2 * time.Hour)).
		WithTagList([]uint64{}).
		WithID(1).
		Json()

	return entryDto, err
}

func (db *decodeBenchmarker) benchDecode_EncodingJson(n int) func(b *testing.B) {
	entry, err := db.getDummyEntryReq()

	if err != nil {
		panic(err)
	}

	return func(b *testing.B) {
		for j := 0; j < n; j++ {
			for i := 0; i < b.N; i++ {
				var req dto.ReqCreateUpdateEntry
				startTime := time.Now()
				err := json.Unmarshal(entry, &req)
				if err != nil {
					panic(err)
				}
				elapsed := time.Since(startTime)
				db.benchMetrics.DecodeRecordTime(encodingJsonLblValue, float64(elapsed.Nanoseconds()))
			}
		}
	}
}

func (db *decodeBenchmarker) benchDecode_EasyJson(n int) func(b *testing.B) {
	entry, err := db.getDummyEntryReq()

	if err != nil {
		panic(err)
	}

	return func(b *testing.B) {
		for j := 0; j < n; j++ {
			for i := 0; i < b.N; i++ {
				var req dto.ReqCreateUpdateEntry
				startTime := time.Now()
				err := easyjson.Unmarshal(entry, &req)
				if err != nil {
					panic(err)
				}
				elapsed := time.Since(startTime)
				db.benchMetrics.DecodeRecordTime(goJsonLblValue, float64(elapsed.Nanoseconds()))
			}
		}
	}
}

func (db *decodeBenchmarker) DecodeBenchmark(n int) (res []string, err error) {
	decodeEncodingJson := db.benchDecode_EncodingJson(n)
	decodeGoJson := db.benchDecode_EasyJson(n)

	resultsEncodingJson := testing.Benchmark(decodeEncodingJson)
	res = append(res, fmt.Sprintf("encoding/json.Unmarshall -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op",
		resultsEncodingJson.N, resultsEncodingJson.NsPerOp(), resultsEncodingJson.AllocsPerOp(), resultsEncodingJson.AllocedBytesPerOp(),
	))

	resultsGoJson := testing.Benchmark(decodeGoJson)
	res = append(res, fmt.Sprintf("easy-json.Unmarshall       -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsGoJson.N, resultsGoJson.NsPerOp(), resultsGoJson.AllocsPerOp(), resultsGoJson.AllocedBytesPerOp(),
	))

	return res, nil
}
