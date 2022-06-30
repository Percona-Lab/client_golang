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
			"collector_scrape_time_ms",
			"Time taken for scrape by collector",
			[]string{"exporter", "collector"},
			nil),
		cache: make(chan Metric, metricsBufferSize),
	}
}

func (m *MetaMetrics) Add(metric Metric) {
	m.cache <- metric
}

func (m *MetaMetrics) Describe(ch chan<- *Desc) {
	ch <- m.desc
}

func (m *MetaMetrics) Collect(ch chan<- Metric) {
	for {
		select {
		case metric := <-m.cache:
			ch <- metric
		default:
			break
		}
	}
}
