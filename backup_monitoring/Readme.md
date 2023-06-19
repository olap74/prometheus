## Prometheus exporter for  S3 based backups monitoring

### General information

This is a Prometheus exporter that can be useful for monitoring backup file size and backup age for backup files located in AWS S3 bucket.
This exporter can be used for any kind of dumps or archives that are being created by some automation. 

It produces four type of metrics: 
 - BackupLatestTimestamp (gauge)
 - BackupOldestTimestamp (gause)
 - BackupLatestSize (gause)
 - BackupOldestSize (gauge)

This allows you to know file size (in bytes) and time stamp (epoch seconds multiplied by 1000) for oldest and latest file in the specified backup directory

This exported is also prepared for using as a docker container and can be executed within EC2, ECS or EKS service. The required permissions (S3 read-only access to specified bucket) can be provided using IaM role. 

This exporter can be configured using environment variables: 
 - `BACKUP_BUCKET` - S3 bucket name
 - `BACKUP_DIR` - path to the directory containing backup files on a specified bucket
 - `EXPORTER_PORT` - Prometheus exporter port. 

Metrics are being produced using the next URL: http(s)://host:port/metrics

### Technical specs

This script has been developed using Golang v.1.14

It uses the next libraries:

    "log"
    "net/http"
    "os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"

Docker image can be created using Dockerfile from the repo. 
The image is being built using multistage build to reduce the Docker image size

