package skizze

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/skizzehq/goskizze/datamodel"
)

func getClient(t *testing.T) (*Client, *fakeSkizze) {
	assert := assert.New(t)

	fs := newFakeSkizze()
	<-fs.ready

	c, err := Dial(fs.address, Options{Insecure: true})
	assert.Nil(err)
	assert.NotNil(c)

	return c, fs
}

func TestDial(t *testing.T) {
	_, fs := getClient(t)
	fs.server.Stop()
}

func TestCreateSnapshot(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer fs.server.Stop()

	rawStatuses := []pb.SnapshotStatus{
		pb.SnapshotStatus_PENDING,
		pb.SnapshotStatus_IN_PROGRESS,
		pb.SnapshotStatus_SUCCESSFUL,
		pb.SnapshotStatus_FAILED,
	}

	statuses := []SnapshotState{
		Pending,
		InProgress,
		Successful,
		Failed,
	}

	for i, rawStatus := range rawStatuses {
		fs.nextReply = &pb.CreateSnapshotReply{
			Status:        &rawStatus,
			StatusMessage: nil,
		}
		fs.nextError = nil

		s, err := c.CreateSnapshot()
		assert.NotNil(s)
		assert.Nil(err)
		assert.Equal(statuses[i], s.Status)
	}
}

func TestGetSnapshot(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer fs.server.Stop()

	rawStatuses := []pb.SnapshotStatus{
		pb.SnapshotStatus_PENDING,
		pb.SnapshotStatus_IN_PROGRESS,
		pb.SnapshotStatus_SUCCESSFUL,
		pb.SnapshotStatus_FAILED,
	}

	statuses := []SnapshotState{
		Pending,
		InProgress,
		Successful,
		Failed,
	}

	for i, rawStatus := range rawStatuses {
		fs.nextReply = &pb.GetSnapshotReply{
			Status:        &rawStatus,
			StatusMessage: nil,
		}
		fs.nextError = nil

		s, err := c.GetSnapshot()
		assert.NotNil(s)
		assert.Nil(err)
		assert.Equal(statuses[i], s.Status)
	}
}
