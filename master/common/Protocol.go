package common

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
