package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"fmt"
	"os"
	"path/filepath"
	"github.com/prometheus/common/log"
)


type dirSizeCol struct {
	dir string
	period time.Duration

	desc *prometheus.Desc

	readCh chan int64
	doneCh chan struct{}
}


func NewDirSizeCol(dir string, period time.Duration) *dirSizeCol {
	c := &dirSizeCol{
		dir: dir,
		period: period,

		desc: prometheus.NewDesc("disk_size_usage",
			"Show how much docker-registry takes disk usage",
			nil, nil,
		),
		readCh: make(chan int64, 1),
		doneCh: make(chan struct{}),
	}

	go c.worker()

	return c
}

func (c *dirSizeCol) Stop() {
	select {
	case <-c.doneCh:
		return
	default:
		close(c.doneCh)
	}
}

// PROMETHEUS API
func (c *dirSizeCol) Describe(ch chan<- *prometheus.Desc) {
	// describe metric
	ch <- c.desc
}

func (c *dirSizeCol) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	//metricValue :=

	//Write latest value for each metric in the prometheus metric channel.
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, c.value + 100)
}

func (c *dirSizeCol) worker() {
	ticker := time.NewTicker(c.period)
	defer ticker.Stop()

	var err error
	var value int64

	for {
		select {
		case c.readCh <- value:
		case <-ticker.C:
			value, err = dirSize(c.dir)
			if err != nil {
				log.Error(err)
			}
		case <-c.doneCh:
			return
		}
	}
}

func dirSize(path string) (int64, error){
	var size int64
	adjSize := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	}

	err := filepath.Walk(path, adjSize)

	return size, err
}