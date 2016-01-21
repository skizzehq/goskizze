package skizze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/skizzehq/goskizze/datamodel"
	. "github.com/skizzehq/goskizze/skizze"
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

func closeAll(c *Client, fs *fakeSkizze) {
	c.Close()
	fs.server.Stop()
}

func stringp(s string) *string {
	return &s
}

func TestDial(t *testing.T) {
	_, fs := getClient(t)
	fs.server.Stop()
}

func TestCreateSnapshot(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

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
	defer closeAll(c, fs)

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

func TestListAll(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	types := []pb.SketchType{pb.SketchType_MEMB, pb.SketchType_FREQ, pb.SketchType_RANK, pb.SketchType_CARD}
	ret := []*pb.Sketch{
		&pb.Sketch{Name: stringp("foobar"), Type: &types[0]},
		&pb.Sketch{Name: stringp("hoobar"), Type: &types[1]},
		&pb.Sketch{Name: stringp("joobar"), Type: &types[2]},
		&pb.Sketch{Name: stringp("loobar"), Type: &types[3]},
	}
	fs.nextReply = &pb.ListReply{
		Sketches: ret,
	}

	sketches, err := c.ListAll()
	assert.NotNil(sketches)
	assert.Nil(err)
	assert.Equal(4, len(sketches))

	for i, sketch := range sketches {
		assert.Equal(ret[i].GetName(), sketch.Name)
	}
}

func TestListSketches(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	types := []pb.SketchType{pb.SketchType_MEMB, pb.SketchType_FREQ, pb.SketchType_RANK, pb.SketchType_CARD}
	stypes := []SketchType{Membership, Frequency, Ranking, Cardinality}
	ret := []*pb.Sketch{
		&pb.Sketch{Name: stringp("foobar"), Type: &types[0]},
		&pb.Sketch{Name: stringp("hoobar"), Type: &types[1]},
		&pb.Sketch{Name: stringp("joobar"), Type: &types[2]},
		&pb.Sketch{Name: stringp("loobar"), Type: &types[3]},
	}
	fs.nextReply = &pb.ListReply{
		Sketches: ret,
	}

	sketches, err := c.ListSketches(Membership)
	assert.NotNil(sketches)
	assert.Nil(err)
	assert.Equal(4, len(sketches))

	req := fs.lastRequest.(*pb.ListRequest)
	assert.Equal(pb.SketchType_MEMB, req.GetType())

	for i, sketch := range sketches {
		assert.Equal(ret[i].GetName(), sketch.Name)
		assert.Equal(stypes[i], sketch.Type)
	}
}

func TestListDomains(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	names := []string{"dom1", "dom2", "dom3", "dom4"}
	fs.nextReply = &pb.ListDomainsReply{Name: names}

	domains, err := c.ListDomains()
	assert.NotNil(domains)
	assert.Nil(err)
	assert.Equal(len(names), len(domains))
	for i, n := range names {
		assert.Equal(n, domains[i])
	}
}

func TestGetDefaults(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	d := Defaults{Rank: 22, Capacity: 10000000}
	fs.nextReply = &pb.Defaults{Rank: &d.Rank, Capacity: &d.Capacity}

	defaults, err := c.GetDefaults()
	assert.Nil(err)
	assert.NotNil(defaults)
	assert.Equal(d.Rank, defaults.Rank)
	assert.Equal(d.Capacity, defaults.Capacity)
}

func TestSetDefaults(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	d := Defaults{Rank: 22, Capacity: 10000000}
	fs.nextReply = &pb.Defaults{Rank: &d.Rank, Capacity: &d.Capacity}

	defaults, err := c.SetDefaults(&d)
	assert.Nil(err)
	assert.NotNil(defaults)
	assert.Equal(d.Rank, defaults.Rank)
	assert.Equal(d.Capacity, defaults.Capacity)

	req := fs.lastRequest.(*pb.Defaults)
	assert.Equal(d.Rank, req.GetRank())
	assert.Equal(d.Capacity, req.GetCapacity())
}

func TestCreateDomain(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Domain{Name: stringp("mydomain")}

	d, err := c.CreateDomain("mydomain", &Defaults{Rank: 1000, Capacity: 100000})
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mydomain", d.Name)

	req := fs.lastRequest.(*pb.Domain)
	assert.Equal("mydomain", req.GetName())
	assert.Equal(int64(1000), req.GetDefaults().GetRank())
	assert.Equal(int64(100000), req.GetDefaults().GetCapacity())
}

func TestDeleteDomain(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Empty{}

	err := c.DeleteDomain("mydomain")
	assert.Nil(err)

	req := fs.lastRequest.(*pb.Domain)
	assert.Equal("mydomain", req.GetName())
}

func TestGetDomain(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Domain{Name: stringp("mydomainman")}

	d, err := c.GetDomain("mydomainman")
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mydomainman", d.Name)

	req := fs.lastRequest.(*pb.Domain)
	assert.Equal("mydomainman", req.GetName())
}
