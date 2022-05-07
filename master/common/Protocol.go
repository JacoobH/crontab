package common

import "encoding/json"

// Job timed task
type Job struct {
	Name     string `json:"name"`     // job name
	Command  string `json:"command"`  // shell command
	CronExpr string `json:"cronExpr"` // cron Expressions
}

// Response HTTP interface response
type Response struct {
	ErrNo int         `json:"errNo"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errNo int, msg string, data interface{}) (resp []byte, err error) {
	// 1. Define a response
	var (
		response Response
	)
	response.ErrNo = errNo
	response.Msg = msg
	response.Data = data

	// 2. Serialize json
	resp, err = json.Marshal(response)
	return
}
