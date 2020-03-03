package main

import (
    "github.com/prometheus/client_golang/prometheus"
)


type filesystemCollector struct {
    sizeMetric *prometheus.Desc
}


func NewFilesystemCollector() *filesystemCollector {
    return &filesystemCollector{
	sizeMetric: prometheus.NewDesc("disk_size_usage",
	    "Show how much docker-registry takes disk usage",
	    nil, nil,
	),
    }
}

func (collector *filesystemCollector) Describe(ch chan<- *prometheus.Desc) {
    // describe metric
    ch <- collector.sizeMetric
}

func (collector *filesystemCollector) Collect(ch chan<- prometheus.Metric) {

    //Implement logic here to determine proper metric value to return to prometheus
    //for each descriptor or call other functions that do so.
    //metricValue :=

    //Write latest value for each metric in the prometheus metric channel.
    ch <- prometheus.MustNewConstMetric(collector.sizeMetric, prometheus.GaugeValue, metricValue + 100)

}
