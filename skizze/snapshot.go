package skizze

import (
	"log"

	pb "github.com/skizzehq/goskizze/datamodel"
)

// SnapshotState indicates the state of the current or previous snapshot.
type SnapshotState int

const (
	// Pending indicates a snapshot is waiting to be taken.
	Pending SnapshotState = iota
	// InProgress indicates a snapshot is currently being taken.
	InProgress
	// Successful indicates the last snapshot was successful.
	Successful
	// Failed indicates the last snapshot had an error.
	Failed
)

// Snapshot represents details of the a snapshot
type Snapshot struct {
	Status  SnapshotState
	Message string
}

func snapshotStatusFromRaw(s pb.SnapshotStatus) SnapshotState {
	switch s {
	case pb.SnapshotStatus_PENDING:
		return Pending
	case pb.SnapshotStatus_IN_PROGRESS:
		return InProgress
	case pb.SnapshotStatus_SUCCESSFUL:
		return Successful
	case pb.SnapshotStatus_FAILED:
		return Failed
	default:
		log.Panicf("Snapshot status %v unknown", s)
	}
	return Pending
}
