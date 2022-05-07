package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacoobH/crontab/master/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// JobMgr job manager
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// G_jobMgr singleton
var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
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

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// SaveJob Save job
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// Save job to /cron/jobs/job_name -> json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)

	// etcd save key
	jobKey = common.JOB_SAVE_DIR + job.Name
	fmt.Println(jobKey)

	// job information json
	if jobValue, err = json.Marshal(*job); err != nil {
		return
	}

	// save to etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	// if it was update, then return old value
	if putResp.PrevKv != nil {
		// deserialize the old values
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// DeleteJob Delete job
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	jobKey = common.JOB_SAVE_DIR + name

	// delete it from etcd
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// return oldJob that be deleted
	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// ListJob List job
func (jobMgr *JobMgr) ListJob() (jobList []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	dirKey = common.JOB_SAVE_DIR

	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// Initialize the array space
	jobList = make([]*common.Job, 0)

	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}
