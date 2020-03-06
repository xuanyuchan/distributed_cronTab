package util

import (
	"context"
	"distributed_cronTab/common"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type Lock struct {
	cli        *clientv3.Client
	key        string
	lease      clientv3.Lease
	leaseId    clientv3.LeaseID
	cancelFunc context.CancelFunc
}

func InitLock(key string) (*Lock, error) {
	clientCfg := clientv3.Config{
		Endpoints:   []string{"119.28.150.217:2379"},
		DialTimeout: 3000 * time.Millisecond,
	}
	cli, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, err
	}
	return &Lock{
		cli:   cli,
		lease: clientv3.NewLease(cli),
		key:   key,
	}, nil
}

func (l *Lock) Lock() error {
	leaseResp, err := l.lease.Grant(context.Background(), 5)
	if err != nil {
		log.Printf("etcd distributed lock, Lock error:%v\n", err)
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	keepRespCh, err := l.lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		log.Printf("etcd distributed lock, Lock keepalive error: %v\n", err)
		cancelFunc()
		l.lease.Revoke(context.Background(), leaseResp.ID)
		return err
	}

	l.leaseId = leaseResp.ID
	l.cancelFunc = cancelFunc

	go func() {
		for {
			select {
			case keepResp := <-keepRespCh:
				if keepResp == nil {
					break
				}
			}
		}
	}()

	//log.Printf("lock revision: %d\n", lastRevision)
	txn := clientv3.NewKV(l.cli).Txn(context.Background())
	lockKey := common.LockDir + l.key
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseResp.ID))).
		Else(clientv3.OpGet(lockKey))

	txnResp, err := txn.Commit()
	if err != nil {
		log.Printf("lock txn commit error:%v\n", err)
		return err
	}
	if !txnResp.Succeeded {
		return errors.New("etcd distributed lock, lock fail")
	}
	return nil
}

func (l *Lock) UnLock() {
	l.lease.Revoke(context.Background(), l.leaseId)
	l.cancelFunc()
}
