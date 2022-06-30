// Copyright 2022 Percona LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"log"
)

const (
	metricsBufferSize = 1024 * 8
)

var (
	MetricsCollector = NewMetaMetricsCollector()
)

type MetaMetrics struct {
	desc  *Desc
	cache chan Metric
}

func NewMetaMetricsCollector() *MetaMetrics {
	return &MetaMetrics{
		desc: NewDesc(
			"collector_meta",
			"Meta Collector",
			[]string{},
			nil),
		cache: make(chan Metric, metricsBufferSize),
	}
}

// Add adds metric to cache. If cache is full, metric is skipped not to block process.
func (m *MetaMetrics) Add(metric Metric) {
	select {
	case m.cache <- metric:
	default: // if buffer is full, discard metric
		//TODO: is there a better way to log?
		log.Println("MetaMetrics blocked")
	}

}

func (m *MetaMetrics) Describe(ch chan<- *Desc) {
	ch <- m.desc
}

// Collect collected buffered metrics, after collecting all of them we exit.
func (m *MetaMetrics) Collect(ch chan<- Metric) {
collect:
	for {
		select {
		case metric := <-m.cache:
			ch <- metric
		default: // exit after we collected all metrics
			break collect
		}
	}
}
