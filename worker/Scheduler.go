package worker

import "github.com/JacoobH/crontab/common"

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

// scheduling coroutine
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent *common.JobEvent
	)
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
