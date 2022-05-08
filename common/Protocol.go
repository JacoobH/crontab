package common

import (
	"encoding/json"
	"strings"
)

// Job timed task
type Job struct {
	Name     string `json:"name" form:"name"`         // job name
	Command  string `json:"command" form:"command"`   // shell command
	CronExpr string `json:"cronExpr" form:"cronExpr"` // cron Expressions
}

// Response HTTP interface response
type Response struct {
	ErrNo int         `json:"errNo"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// JobEvent Change event
type JobEvent struct {
	EventType int //SAVE | DELETE
	Job       *Job
}

func BuildResponse(errNo int, msg string, data interface{}) (resp Response) {
	// 1. Define a response
	var (
		response Response
	)
	response.ErrNo = errNo
	response.Msg = msg
	response.Data = data

	// 2. Serialize json
	resp = response
	return
}

func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// ExtractJobName Extract the task name from the key of etCD
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}
