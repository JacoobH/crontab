package master

import (
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
