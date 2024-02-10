package metrics

import (
	"strings"
)

func NewCounter(name string) MetricsCounter {
	newCounter := MetricsCounter{name: name, tags: [][2]string{}}
	return newCounter
}

type MetricsCounter struct {
	name         string
	tags         [][2]string
	IncreaseFunc func(label string, offset uint64)
}

func (counter *MetricsCounter) AddTag(tag string, value string) *MetricsCounter {
	counter.tags = append(counter.tags, [2]string{tag, value})
	return counter
}

func (counter *MetricsCounter) Increase(offset uint64) {
	kvPairs := []string{}
	for _, tag := range counter.tags {
		key, value := tag[0], tag[1]
		kvPairs = append(kvPairs, key+"="+value)
	}
	counter.tags = [][2]string{}

	if counter.IncreaseFunc != nil {
		counter.IncreaseFunc(counter.name+"_"+strings.Join(kvPairs, "_"), offset)
	}
}
