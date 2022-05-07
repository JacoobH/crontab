package common

// Job timed task
type Job struct {
	Name     string `json:"name"`     // job name
	Command  string `json:"command"`  // shell command
	CronExpr string `json:"cronExpr"` // cron Expressions
}
