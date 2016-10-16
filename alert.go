package main

import (
	"time"

	"github.com/influxdata/influxdb/influxql"
)

type alertData struct {
	ID       string          `json:"id"`
	Message  string          `json:"message"`
	Details  string          `json:"details"`
	Time     time.Time       `json:"time"`
	Duration time.Duration   `json:"duration"`
	Level    string          `json:"level"`
	Data     influxql.Result `json:"data"`
	Status   string          `json:"status"`
}
