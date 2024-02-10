package metrics_test

import (
	"fmt"
	"testing"
	"wasm-jwt-filter/pkg/metrics"
)

var counter = metrics.NewCounter("test_counts")

var value string

func DoIncrease(label string, offset uint64) {
	value = fmt.Sprintf("%s %d", label, offset)
}

func TestCounter(t *testing.T) {
	counter.IncreaseFunc = DoIncrease
	counter.AddTag("tag1", "value1")
	counter.AddTag("tag2", "value2")
	counter.Increase(1)
	if value != "test_counts_tag1=value1_tag2=value2 1" {
		t.Error(nil)
	}

	counter.AddTag("tag3", "value3")
	counter.AddTag("tag4", "value4")
	counter.Increase(1)
	if value != "test_counts_tag3=value3_tag4=value4 1" {
		t.Error(nil)
	}
}
