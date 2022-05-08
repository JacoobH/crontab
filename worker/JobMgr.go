package worker

import (
	"context"
	"fmt"
	"github.com/JacoobH/crontab/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// JobMgr job manager
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

// G_jobMgr singleton
var (
	G_jobMgr *JobMgr
)

// Monitoring job changes
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	//Get the all jobs of /corn/jobs and observe subsequent changes
	if getResp, err = G_jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	}

	// current jobs
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// TODO:把这个job同步给scheduler(调度协程)
			return
		}
	}

	go func() { // Listening coroutines
		//Etcd Cluster transaction ID, monotonically increasing
		watchStartRevision = getResp.Header.Revision + 1
		//Start watching
		watchChan = G_jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision))
		//Handle event about the changes of kv
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					//Building a update Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					//Building a delete Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				//TODO:Push an event to the Scheduler
			}
		}
	}()
	return
}

func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,                                     //cluster network address
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, // timeout
	}

	//Establish a client
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//Use to read or write KV of etcd
	kv = clientv3.NewKV(client)

	//Apply for a lease
	lease = clientv3.NewLease(client)

	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	G_jobMgr.watchJobs()

	return
}
