package common

const (
	// JOB_SAVE_DIR Job save directory
	JOB_SAVE_DIR = "/cron/jobs/"

	// JOB_KILLER_DIR kill directory
	JOB_KILLER_DIR = "/cron/killer/"

	JOB_LOCK_DIR = "/cron/lock/"

	// JOB_EVENT_SAVE save job event
	JOB_EVENT_SAVE = 1

	// JOB_EVENT_DELETE delete job event
	JOB_EVENT_DELETE = 2

	JOB_EVENT_KILL = 3
)
