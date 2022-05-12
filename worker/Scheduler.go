package worker

import (
	"github.com/JacoobH/crontab/common"
	"time"
)

// Scheduler job scheduling
type Scheduler struct {
	jobEventChan chan *common.JobEvent              // Etcd job event queue
	jobPlanTable map[string]*common.JobSchedulePlan // Job scheduling schedule
}

var (
	G_scheduler *Scheduler
)

// Processing job Events
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExisted      bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// TrySchedule Recalculate the task scheduling status
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)

	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}
	// current time
	now = time.Now()
	//1. Iterate through all jobs
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO: try to exec job
			jobPlan.NextTime = jobPlan.Expr.Next(now) // Updated the next execution time
		}
		// Count the last time a job expired
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	//Interval for next scheduling（earTime - now）
	scheduleAfter = (*nearTime).Sub(now)
	return
	//2. Expired jobs are executed immediately
	//3. Count the time of the most recent expired job (N seconds after expiration == scheduleAfter)
}

// scheduling coroutine
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	//Initialize(1sec)
	scheduleAfter = scheduler.TrySchedule()

	//Delay timer for scheduling
	scheduleTimer = time.NewTimer(scheduleAfter)
	// Timing job common.Job
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: // Listen for job change events
			//CRUD the job list maintained in memory
			scheduler.handleJobEvent(jobEvent)
		}
	}
}

// PushJobEvent Pushing job change events
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// InitScheduler Initialize scheduler
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
	}
	// Start scheduling coroutine
	go G_scheduler.scheduleLoop()
	return
}
