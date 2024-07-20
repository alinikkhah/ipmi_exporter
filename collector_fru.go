// Copyright 2021 The Prometheus Authors
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

package main

import (
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus-community/ipmi_exporter/freeipmi"
)

const (
	FRUCollectorName CollectorName = "fru"
)

var (
	fruInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "fru", "info"),
		"Constant metric with value '1' providing details about the FRU.",
		[]string{"serial_number"},
		nil,
	)
)

type FRUCollector struct{}

func (c FRUCollector) Name() CollectorName {
	return FRUCollectorName
}

func (c FRUCollector) Cmd() string {
	return "ipmi-fru"
}

func (c FRUCollector) Args() []string {
	return []string{"--device-id=0"}
}

func (c FRUCollector) Collect(result freeipmi.Result, ch chan<- prometheus.Metric, target ipmiTarget) (int, error) {
	fruProductSerial, err := freeipmi.GetFRUProductSerialNumber(result)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to collect FRU data", "target", targetName(target.host), "error", err)
		return 0, err
	}
	ch <- prometheus.MustNewConstMetric(
		fruInfoDesc,
		prometheus.GaugeValue,
		1,
		fruProductSerial,
	)
	return 1, nil
}
