package replication

import (
	"context"

	"concurrency_go_course/internal/storage/wal"
)

const (
	ReplicaTypeMaster = "master"
	ReplicaTypeSlave  = "slave"
)

type Replication interface {
	Start(context.Context) error
	IsMaster() bool
	ReplicationStream() chan []wal.Request
}
