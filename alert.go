package main

import (
	"time"

	"github.com/influxdata/influxdb/influxql"
)

type AlertData struct {
	ID       string          `json:"id"`
	Message  string          `json:"message"`
	Details  string          `json:"details"`
	Time     time.Time       `json:"time"`
	Duration time.Duration   `json:"duration"`
	Level    string          `json:"level"`
	Data     influxql.Result `json:"data"`
	Assignee string          `json:"assignee"`
	Status   string          `json:"status"`
}
