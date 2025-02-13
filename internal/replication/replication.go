package replication

const (
	// ReplicaTypeMaster is replication type master
	ReplicaTypeMaster = "master"
	// ReplicaTypeSlave is replication type slave
	ReplicaTypeSlave = "slave"
)

// Replication is struct for replication
type Replication struct {
	Slave  *Slave
	Master *Master
}
