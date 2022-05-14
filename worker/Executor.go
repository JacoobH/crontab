package worker

import (
	"context"
	"github.com/JacoobH/crontab/common"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(jobExecuteInfo *common.JobExecuteInfo) {
	go func() {
		var (
			cmd              *exec.Cmd
			err              error
			output           []byte
			jobExecuteResult *common.JobExecuteResult
		)
		jobExecuteResult = &common.JobExecuteResult{
			JobExecuteInfo: jobExecuteInfo,
			OutPut:         make([]byte, 0),
		}
		jobExecuteResult.StartTime = time.Now()

		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", jobExecuteInfo.Job.Command)
		output, err = cmd.CombinedOutput()
		jobExecuteResult.EndTime = time.Now()
		jobExecuteResult.OutPut = output
		jobExecuteResult.Err = err
		// When the job is completed, the result of the execution is returned to the Scheduler and deletes the execution record from the executingTable
		G_scheduler.PushJobResult(jobExecuteResult)
	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
