package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// Job timed job
type Job struct {
	Name     string `json:"name" form:"name"`         // job name
	Command  string `json:"command" form:"command"`   // shell command
	CronExpr string `json:"cronExpr" form:"cronExpr"` // cron Expressions
}

// Job scheduling plan
type JobSchedulePlan struct {
	Job      *Job                 // Information about the tasks to be scheduled
	Expr     *cronexpr.Expression // The resolved cronexpr expression
	NextTime time.Time            // Next scheduling time
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
	job = &Job{}
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

// BuildJobSchedulePlan Construct a job execution plan
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	// Parse the cron expression of the job
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {

	}

	// Generate JobSchedulePlan object
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}
