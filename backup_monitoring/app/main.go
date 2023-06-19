package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	backupBucket    = os.Getenv("BACKUP_BUCKET")
	backupDirectory = os.Getenv("BACKUP_DIR")
	exporterPort    = os.Getenv("EXPORTER_PORT")
)

type s3Collector struct {
	BackupLatestTimestamp *prometheus.Desc
	BackupLatestSize      *prometheus.Desc
	BackupOldestTimestamp *prometheus.Desc
	BackupOldestSize      *prometheus.Desc
}

func news3Collector() *s3Collector {
	var (
		bucketLabel = prometheus.Labels{"backup_bucket": backupBucket}
	)

	return &s3Collector{
		BackupLatestTimestamp: prometheus.NewDesc("backup_latest_file_timestamp",
			"Last modified timestamp(milliseconds) for latest backup file",
			nil, bucketLabel,
		),
		BackupLatestSize: prometheus.NewDesc("backup_latest_file_size",
			"Size in bytes for latest backup file",
			nil, bucketLabel,
		),
		BackupOldestTimestamp: prometheus.NewDesc("backup_oldest_file_timestamp",
			"Last modified timestamp(milliseconds) for oldest backup file",
			nil, bucketLabel,
		),
		BackupOldestSize: prometheus.NewDesc("backup_oldest_file_size",
			"Size in bytes for oldest backup file",
			nil, bucketLabel,
		),
	}
}

func (collector *s3Collector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.BackupLatestTimestamp
	ch <- collector.BackupLatestSize
	ch <- collector.BackupOldestTimestamp
	ch <- collector.BackupOldestSize
}

func (collector *s3Collector) Collect(ch chan<- prometheus.Metric) {

	var LatestTimestampValue int64 = 0
	var LatestSizeValue int64 = 0
	var OldestTimestampValue int64 = 0
	var OldestSizeValue int64 = 0

	svc := s3.New(session.New())
	params := &s3.ListObjectsInput{
		Bucket: aws.String(backupBucket),
		Prefix: aws.String(backupDirectory),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Printf("AWS connection error: %s\n", err)
	}
	data := resp.Contents

	for _, key := range data {
		if *key.Key != backupDirectory {
			t := *key.LastModified
			var timestamp int64 = t.Unix() * 1000
			if timestamp > LatestTimestampValue {
				LatestTimestampValue = timestamp
				LatestSizeValue = *key.Size
			}
			if OldestTimestampValue == 0 {
				OldestTimestampValue = timestamp
				OldestSizeValue = *key.Size
			}
			if timestamp < OldestTimestampValue {
				OldestTimestampValue = timestamp
				OldestSizeValue = *key.Size
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(collector.BackupLatestTimestamp, prometheus.GaugeValue, float64(LatestTimestampValue))
	ch <- prometheus.MustNewConstMetric(collector.BackupLatestSize, prometheus.GaugeValue, float64(LatestSizeValue))
	ch <- prometheus.MustNewConstMetric(collector.BackupOldestTimestamp, prometheus.GaugeValue, float64(OldestTimestampValue))
	ch <- prometheus.MustNewConstMetric(collector.BackupOldestSize, prometheus.GaugeValue, float64(OldestSizeValue))
}

func main() {
	metrics := news3Collector()
	prometheus.MustRegister(metrics)

	log.Printf("Starting web server at %s\n", exporterPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+exporterPort, nil))
}
