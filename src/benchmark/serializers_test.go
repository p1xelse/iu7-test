package benchmark

import (
	"fmt"
	"testing"
	"timetracker/benchmark/mocks"

	"github.com/golang/mock/gomock"
)

const benchN = 1

func TestEncodeDecode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metricsMock := mocks.NewMockMetrics(ctrl)
	metricsMock.EXPECT().DecodeRecordTime(gomock.Any(), gomock.Any()).AnyTimes()
	metricsMock.EXPECT().EncodeRecordTime(gomock.Any(), gomock.Any()).AnyTimes()

	encodeBenchmarker := NewEncodeBenchmarker(metricsMock)
	decodeBenchmarker := NewDecodeBenchmarker(metricsMock)

	resEncode, err := encodeBenchmarker.EncodeBenchmark(benchN)
	if err != nil {
		panic(err)
	}
	fmt.Println("Encode result: ")
	for _, res := range resEncode {
		fmt.Println(res)
	}

	resDecode, err := decodeBenchmarker.DecodeBenchmark(benchN)
	if err != nil {
		panic(err)
	}
	fmt.Println("Decode result: ")
	for _, res := range resDecode {
		fmt.Println(res)
	}
}
