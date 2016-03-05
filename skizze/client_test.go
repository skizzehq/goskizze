package skizze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/skizzehq/goskizze/protobuf"
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
	fs.nextReply = &pb.ListDomainsReply{Names: names}

	domains, err := c.ListDomains()
	assert.NotNil(domains)
	assert.Nil(err)
	assert.Equal(len(names), len(domains))
	for i, n := range names {
		assert.Equal(n, domains[i])
	}
}

func TestCreateDomain(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Domain{Name: stringp("mydomain")}

	d, err := c.CreateDomain("mydomain")
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mydomain", d.Name)

	req := fs.lastRequest.(*pb.Domain)
	assert.Equal("mydomain", req.GetName())
}

func TestCreateDomainWithProperties(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Domain{
		Name: stringp("mydomain"),
	}

	d, err := c.CreateDomainWithProperties("mydomain", &DomainProperties{
		MembershipProperties: Properties{
			MaxUniqueItems: 256,
			ErrorRate:      0.6,
		},
		FrequencyProperties: Properties{
			MaxUniqueItems: 512,
			ErrorRate:      0.8,
		},
		RankingsProperties: Properties{
			Size: 101,
		},
	})
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mydomain", d.Name)

	req := fs.lastRequest.(*pb.Domain)
	assert.Equal("mydomain", req.GetName())
	assert.Equal(int64(256), req.GetSketches()[0].GetProperties().GetMaxUniqueItems())
	assert.Equal(float32(0.6), req.GetSketches()[0].GetProperties().GetErrorRate())
	assert.Equal(int64(512), req.GetSketches()[1].GetProperties().GetMaxUniqueItems())
	assert.Equal(float32(0.8), req.GetSketches()[1].GetProperties().GetErrorRate())
	assert.Equal(int64(101), req.GetSketches()[2].GetProperties().GetSize())
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

func TestCreateSketch(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	st := pb.SketchType_RANK
	fs.nextReply = &pb.Sketch{Name: stringp("mysketch"), Type: &st}

	d, err := c.CreateSketch("mysketch", Ranking, &Properties{Size: 1000})
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mysketch", d.Name)
	assert.Equal(Ranking, d.Type)

	req := fs.lastRequest.(*pb.Sketch)
	assert.Equal("mysketch", req.GetName())
	assert.Equal(int64(1000), req.GetProperties().GetSize())
}

func TestDeleteSketch(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.Empty{}

	err := c.DeleteSketch("mysketch", Frequency)
	assert.Nil(err)

	req := fs.lastRequest.(*pb.Sketch)
	assert.Equal("mysketch", req.GetName())
}

func TestGetSketch(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	st := pb.SketchType_CARD
	fs.nextReply = &pb.Sketch{Name: stringp("mysketchman"), Type: &st}

	d, err := c.GetSketch("mysketchman", Cardinality)
	assert.Nil(err)
	assert.NotNil(d)
	assert.Equal("mysketchman", d.Name)
	assert.Equal(Cardinality, d.Type)

	req := fs.lastRequest.(*pb.Sketch)
	assert.Equal("mysketchman", req.GetName())
}

func TestAddToSketch(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.AddReply{}

	values := []string{"one", "two", "three", "four", "five"}
	err := c.AddToSketch("mysketchysketch", Cardinality, values...)
	assert.Nil(err)

	req := fs.lastRequest.(*pb.AddRequest)
	assert.Nil(req.Domain)
	assert.NotNil(req.Sketch)
	assert.Equal(pb.SketchType_CARD, req.GetSketch().GetType())
	assert.Equal(len(values), len(req.GetValues()))
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestAddToDomain(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	fs.nextReply = &pb.AddReply{}

	values := []string{"one", "two", "three", "four", "five", "six"}
	err := c.AddToDomain("mysketchysketch", values...)
	assert.Nil(err)

	req := fs.lastRequest.(*pb.AddRequest)
	assert.Nil(req.Sketch)
	assert.NotNil(req.Domain)
	assert.Equal(len(values), len(req.GetValues()))
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestGetMembership(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	yes := true
	no := false
	r := &pb.MembershipResult{}
	r.Memberships = []*pb.Membership{
		&pb.Membership{Value: stringp("one"), IsMember: &yes},
		&pb.Membership{Value: stringp("two"), IsMember: &no},
		&pb.Membership{Value: stringp("one"), IsMember: &yes},
	}
	fs.nextReply = &pb.GetMembershipReply{Results: []*pb.MembershipResult{r}}

	values := []string{"foo", "bar", "baz"}
	m, err := c.GetMembership("mymembers", values...)
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(len(values), len(m))
	for i, v := range r.Memberships {
		assert.Equal(v.GetValue(), m[i].Value)
		assert.Equal(v.GetIsMember(), m[i].IsMember)
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal(pb.SketchType_MEMB, req.GetSketches()[0].GetType())
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestGetMultiMembership(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	yes := true
	no := false
	r := &pb.MembershipResult{}
	r.Memberships = []*pb.Membership{
		&pb.Membership{Value: stringp("one"), IsMember: &yes},
		&pb.Membership{Value: stringp("two"), IsMember: &no},
		&pb.Membership{Value: stringp("one"), IsMember: &yes},
	}
	fs.nextReply = &pb.GetMembershipReply{Results: []*pb.MembershipResult{r, r, r}}

	values := []string{"foo", "bar", "baz"}
	m, err := c.GetMultiMembership([]string{"mymembers", "myothermembers"}, values...)
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(3, len(m))
	for _, result := range m {
		for i, v := range r.Memberships {
			assert.Equal(v.GetValue(), result[i].Value)
			assert.Equal(v.GetIsMember(), result[i].IsMember)
		}
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal("myothermembers", req.GetSketches()[1].GetName())
	assert.Equal(pb.SketchType_MEMB, req.GetSketches()[0].GetType())
	assert.Equal(pb.SketchType_MEMB, req.GetSketches()[1].GetType())
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestGetFrequency(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	one := int64(1)
	thou := int64(1000)

	r := &pb.FrequencyResult{}
	r.Frequencies = []*pb.Frequency{
		&pb.Frequency{Value: stringp("one"), Count: &one},
		&pb.Frequency{Value: stringp("two"), Count: &thou},
		&pb.Frequency{Value: stringp("one"), Count: &one},
	}
	fs.nextReply = &pb.GetFrequencyReply{Results: []*pb.FrequencyResult{r}}

	values := []string{"foo", "bar", "baz"}
	m, err := c.GetFrequency("mymembers", values...)
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(len(values), len(m))
	for i, v := range r.Frequencies {
		assert.Equal(v.GetValue(), m[i].Value)
		assert.Equal(v.GetCount(), m[i].Count)
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal(pb.SketchType_FREQ, req.GetSketches()[0].GetType())
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestGetMultiFrequency(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	one := int64(1)
	thou := int64(1000)
	r := &pb.FrequencyResult{}
	r.Frequencies = []*pb.Frequency{
		&pb.Frequency{Value: stringp("one"), Count: &one},
		&pb.Frequency{Value: stringp("two"), Count: &thou},
		&pb.Frequency{Value: stringp("one"), Count: &one},
	}

	fs.nextReply = &pb.GetFrequencyReply{Results: []*pb.FrequencyResult{r, r, r}}

	values := []string{"foo", "bar", "baz"}
	m, err := c.GetMultiFrequency([]string{"mymembers", "myothermembers"}, values...)
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(3, len(m))
	for _, result := range m {
		for i, v := range r.Frequencies {
			assert.Equal(v.GetValue(), result[i].Value)
			assert.Equal(v.GetCount(), result[i].Count)
		}
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal("myothermembers", req.GetSketches()[1].GetName())
	assert.Equal(pb.SketchType_FREQ, req.GetSketches()[0].GetType())
	assert.Equal(pb.SketchType_FREQ, req.GetSketches()[1].GetType())
	for i, v := range req.GetValues() {
		assert.Equal(values[i], v)
	}
}

func TestGetRankings(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	one := int64(1)
	thou := int64(1000)

	r := &pb.RankingsResult{}
	r.Rankings = []*pb.Rank{
		&pb.Rank{Value: stringp("one"), Count: &one},
		&pb.Rank{Value: stringp("two"), Count: &thou},
		&pb.Rank{Value: stringp("one"), Count: &one},
	}
	fs.nextReply = &pb.GetRankingsReply{Results: []*pb.RankingsResult{r}}

	m, err := c.GetRankings("mymembers")
	assert.Nil(err)
	assert.NotNil(m)
	for i, v := range r.Rankings {
		assert.Equal(v.GetValue(), m[i].Value)
		assert.Equal(v.GetCount(), m[i].Count)
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal(pb.SketchType_RANK, req.GetSketches()[0].GetType())
}

func TestGetMultiRankings(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	one := int64(1)
	thou := int64(1000)
	r := &pb.RankingsResult{}
	r.Rankings = []*pb.Rank{
		&pb.Rank{Value: stringp("one"), Count: &one},
		&pb.Rank{Value: stringp("two"), Count: &thou},
		&pb.Rank{Value: stringp("one"), Count: &one},
	}
	fs.nextReply = &pb.GetRankingsReply{Results: []*pb.RankingsResult{r, r, r}}

	m, err := c.GetMultiRankings([]string{"mymembers", "myothermembers"})
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(3, len(m))
	for _, result := range m {
		for i, v := range r.Rankings {
			assert.Equal(v.GetValue(), result[i].Value)
			assert.Equal(v.GetCount(), result[i].Count)
		}
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal("myothermembers", req.GetSketches()[1].GetName())
	assert.Equal(pb.SketchType_RANK, req.GetSketches()[0].GetType())
	assert.Equal(pb.SketchType_RANK, req.GetSketches()[1].GetType())
}

func TestGetCardinality(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	thou := int64(1000)
	fs.nextReply = &pb.GetCardinalityReply{
		Results: []*pb.CardinalityResult{
			&pb.CardinalityResult{
				Cardinality: &thou,
			},
		},
	}

	card, err := c.GetCardinality("mymembers")
	assert.Nil(err)
	assert.Equal(thou, card)

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal(pb.SketchType_CARD, req.GetSketches()[0].GetType())
}

func TestGetMultiCardinality(t *testing.T) {
	assert := assert.New(t)

	c, fs := getClient(t)
	defer closeAll(c, fs)

	thou := int64(1000)
	fs.nextReply = &pb.GetCardinalityReply{
		Results: []*pb.CardinalityResult{
			&pb.CardinalityResult{
				Cardinality: &thou,
			},
			&pb.CardinalityResult{
				Cardinality: &thou,
			},
			&pb.CardinalityResult{
				Cardinality: &thou,
			},
		},
	}

	m, err := c.GetMultiCardinality([]string{"mymembers", "myothermembers"})
	assert.Nil(err)
	assert.NotNil(m)
	assert.Equal(3, len(m))
	for _, result := range m {
		assert.Equal(result, thou)
	}

	req := fs.lastRequest.(*pb.GetRequest)
	assert.Equal("mymembers", req.GetSketches()[0].GetName())
	assert.Equal("myothermembers", req.GetSketches()[1].GetName())
	assert.Equal(pb.SketchType_CARD, req.GetSketches()[0].GetType())
	assert.Equal(pb.SketchType_CARD, req.GetSketches()[1].GetType())
}
