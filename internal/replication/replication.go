package replication

import (
	"context"

	"concurrency_go_course/internal/storage/wal"
)

const (
	// ReplicaTypeMaster is replication type master
	ReplicaTypeMaster = "master"
	// ReplicaTypeSlave is replication type slave
	ReplicaTypeSlave = "slave"
)

// Replication is interface for replication
type Replication interface {
	Start(context.Context)
	IsMaster() bool
	ReplicationStream() chan []wal.Request
}
