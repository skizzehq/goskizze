package skizze

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/skizzehq/goskizze/datamodel"
)

var (
	typeMemb = pb.SketchType_MEMB
	typeFreq = pb.SketchType_FREQ
	typeRank = pb.SketchType_RANK
	typeCard = pb.SketchType_CARD
)

// Client represents a a thread-safe connection to Skizze
type Client struct {
	opts Options

	conn   *grpc.ClientConn
	client pb.SkizzeClient
}

// Dial initalizes a connection to Skizze and returns a client
func Dial(address string, opts Options) (*Client, error) {
	var gOpts []grpc.DialOption
	if opts.Insecure == true {
		gOpts = append(gOpts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, gOpts...)
	if err != nil {
		return nil, fmt.Errorf("Unable to dial Skizze at %v: %v", address, err)
	}

	return &Client{
		opts:   opts,
		conn:   conn,
		client: pb.NewSkizzeClient(conn),
	}, nil
}

// Close shuts down the client connection to Skizze.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// CreateSnapshot queues a snapshot operation.
func (c *Client) CreateSnapshot() (*Snapshot, error) {
	reply, err := c.client.CreateSnapshot(context.Background(), &pb.CreateSnapshotRequest{})
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		Status:  snapshotStatusFromRaw(reply.GetStatus()),
		Message: reply.GetStatusMessage(),
	}, nil
}

// GetSnapshot retrieves the information on the current or last snapshot.
func (c *Client) GetSnapshot() (*Snapshot, error) {
	reply, err := c.client.GetSnapshot(context.Background(), &pb.GetSnapshotRequest{})
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		Status:  snapshotStatusFromRaw(reply.GetStatus()),
		Message: reply.GetStatusMessage(),
	}, nil
}

// ListAll gets all the available Sketches.
func (c *Client) ListAll() (ret []*Sketch, err error) {
	reply, err := c.client.ListAll(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, err
	}
	for _, sketch := range reply.GetSketches() {
		ret = append(ret, newSketchFromRaw(sketch))
	}
	return ret, err
}

// ListSketches gets all the sketches of the specified type.
func (c *Client) ListSketches(t SketchType) (ret []*Sketch, err error) {
	rt := getRawSketchForSketchType(t)
	reply, err := c.client.List(context.Background(), &pb.ListRequest{Type: &rt})
	if err != nil {
		return nil, err
	}
	for _, sketch := range reply.GetSketches() {
		ret = append(ret, newSketchFromRaw(sketch))
	}
	return ret, err
}

// ListDomains gets all the available domains
func (c *Client) ListDomains() (ret []string, err error) {
	reply, err := c.client.ListDomains(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, err
	}
	return reply.GetNames(), nil
}

// CreateDomain creates a new domain with default properties per Sketch.
func (c *Client) CreateDomain(name string) (*Domain, error) {
	rd := &pb.Domain{Name: &name}

	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name: &name,
		Type: &typeMemb,
		Properties: &pb.SketchProperties{
			MaxUniqueItems: &defaultMembUnique,
			ErrorRate:      &defaultMembErrorRate,
		},
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name: &name,
		Type: &typeFreq,
		Properties: &pb.SketchProperties{
			MaxUniqueItems: &defaultFreqUnique,
			ErrorRate:      &defaultFreqErrorRate,
		},
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name: &name,
		Type: &typeRank,
		Properties: &pb.SketchProperties{
			Size: &defaultRankSize,
		},
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name: &name,
		Type: &typeCard,
	})

	reply, err := c.client.CreateDomain(context.Background(), rd)
	if err != nil {
		return nil, err
	}
	return newDomainFromRaw(reply), nil
}

// CreateDomainWithProperties creates a domain with customized properties.
func (c *Client) CreateDomainWithProperties(name string, props *DomainProperties) (*Domain, error) {
	rd := &pb.Domain{Name: &name}

	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name:       &name,
		Type:       &typeMemb,
		Properties: newRawPropertiesFromProperties(&props.MembershipProperties),
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name:       &name,
		Type:       &typeFreq,
		Properties: newRawPropertiesFromProperties(&props.FrequencyProperties),
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name:       &name,
		Type:       &typeRank,
		Properties: newRawPropertiesFromProperties(&props.RankingsProperties),
	})
	rd.Sketches = append(rd.Sketches, &pb.Sketch{
		Name: &name,
		Type: &typeCard,
	})

	reply, err := c.client.CreateDomain(context.Background(), rd)
	if err != nil {
		return nil, err
	}
	return newDomainFromRaw(reply), nil
}

// DeleteDomain deletes a domain
func (c *Client) DeleteDomain(name string) error {
	rd := &pb.Domain{Name: &name}
	_, err := c.client.DeleteDomain(context.Background(), rd)
	if err != nil {
		return err
	}
	return nil
}

// GetDomain gets the details of a domain.
func (c *Client) GetDomain(name string) (*Domain, error) {
	rd := &pb.Domain{Name: &name}
	reply, err := c.client.GetDomain(context.Background(), rd)
	if err != nil {
		return nil, err
	}
	return newDomainFromRaw(reply), nil
}

// CreateSketch creates a new sketch.
func (c *Client) CreateSketch(name string, t SketchType, p *Properties) (*Sketch, error) {
	rt := getRawSketchForSketchType(t)
	rd := &pb.Sketch{Name: &name, Type: &rt, Properties: newRawPropertiesFromProperties(p)}
	reply, err := c.client.CreateSketch(context.Background(), rd)
	if err != nil {
		return nil, err
	}
	return newSketchFromRaw(reply), nil
}

// DeleteSketch deletes a sketch
func (c *Client) DeleteSketch(name string, t SketchType) error {
	rt := getRawSketchForSketchType(t)
	rd := &pb.Sketch{Name: &name, Type: &rt}
	_, err := c.client.DeleteSketch(context.Background(), rd)
	if err != nil {
		return err
	}
	return nil
}

// GetSketch gets the details of a sketch.
func (c *Client) GetSketch(name string, t SketchType) (*Sketch, error) {
	rt := getRawSketchForSketchType(t)
	rd := &pb.Sketch{Name: &name, Type: &rt}
	reply, err := c.client.GetSketch(context.Background(), rd)
	if err != nil {
		return nil, err
	}
	return newSketchFromRaw(reply), nil
}

// AddToSketch will add the supplied values to the sketch's data set.
func (c *Client) AddToSketch(name string, t SketchType, values ...string) error {
	rt := getRawSketchForSketchType(t)
	rs := pb.Sketch{Name: &name, Type: &rt}
	_, err := c.client.Add(context.Background(), &pb.AddRequest{Sketch: &rs, Values: values})
	return err
}

// AddToDomain will add the supplied values to the domain's data set.
func (c *Client) AddToDomain(name string, values ...string) error {
	rd := pb.Domain{Name: &name}
	_, err := c.client.Add(context.Background(), &pb.AddRequest{Domain: &rd, Values: values})
	return err
}

// GetMembership queries the sketch for membership (true/false) for the provided values.
func (c *Client) GetMembership(name string, values ...string) (ret []*MembershipResult, err error) {
	rs := pb.Sketch{Name: &name, Type: &typeMemb}
	reply, err := c.client.GetMembership(context.Background(), &pb.GetRequest{Sketches: []*pb.Sketch{&rs}, Values: values})
	if err != nil {
		return nil, err
	}
	for _, m := range reply.GetResults()[0].GetMemberships() {
		ret = append(ret, &MembershipResult{Value: m.GetValue(), IsMember: m.GetIsMember()})
	}
	return ret, nil
}

// GetMultiMembership queries multiple sketches for membership of the provided values.
func (c *Client) GetMultiMembership(names []string, values ...string) (ret [][]*MembershipResult, err error) {
	req := &pb.GetRequest{Values: values}
	for i := range names {
		req.Sketches = append(req.Sketches, &pb.Sketch{Name: &names[i], Type: &typeMemb})
	}

	reply, err := c.client.GetMembership(context.Background(), req)
	if err != nil {
		return nil, err
	}

	for _, result := range reply.GetResults() {
		r := []*MembershipResult{}
		for _, m := range result.GetMemberships() {
			r = append(r, &MembershipResult{Value: m.GetValue(), IsMember: m.GetIsMember()})
		}
		ret = append(ret, r)
	}
	return ret, nil
}

// GetFrequency queries the sketch for frequency for the provided values.
func (c *Client) GetFrequency(name string, values ...string) (ret []*FrequencyResult, err error) {
	rs := pb.Sketch{Name: &name, Type: &typeFreq}
	reply, err := c.client.GetFrequency(context.Background(), &pb.GetRequest{Sketches: []*pb.Sketch{&rs}, Values: values})
	if err != nil {
		return nil, err
	}
	for _, m := range reply.GetResults()[0].GetFrequencies() {
		ret = append(ret, &FrequencyResult{Value: m.GetValue(), Count: m.GetCount()})
	}
	return ret, nil
}

// GetMultiFrequency queries multiple sketches for the frequency of the provided values.
func (c *Client) GetMultiFrequency(names []string, values ...string) (ret [][]*FrequencyResult, err error) {
	req := &pb.GetRequest{Values: values}
	for i := range names {
		req.Sketches = append(req.Sketches, &pb.Sketch{Name: &names[i], Type: &typeFreq})
	}
	reply, err := c.client.GetFrequency(context.Background(), req)
	if err != nil {
		return nil, err
	}
	for _, result := range reply.GetResults() {
		r := []*FrequencyResult{}
		for _, m := range result.GetFrequencies() {
			r = append(r, &FrequencyResult{Value: m.GetValue(), Count: m.GetCount()})
		}
		ret = append(ret, r)
	}
	return ret, nil
}

// GetRankings queries the sketch for the top rankings.
func (c *Client) GetRankings(name string) (ret []*RankingsResult, err error) {
	rs := pb.Sketch{Name: &name, Type: &typeRank}
	reply, err := c.client.GetRankings(context.Background(), &pb.GetRequest{Sketches: []*pb.Sketch{&rs}})
	if err != nil {
		return nil, err
	}
	for _, m := range reply.GetResults()[0].GetRankings() {
		ret = append(ret, &RankingsResult{Value: m.GetValue(), Count: m.GetCount()})
	}
	return ret, nil
}

// GetMultiRankings queries multiple sketches for the top rankings.
func (c *Client) GetMultiRankings(names []string) (ret [][]*RankingsResult, err error) {
	req := &pb.GetRequest{}
	for i := range names {
		req.Sketches = append(req.Sketches, &pb.Sketch{Name: &names[i], Type: &typeRank})
	}
	reply, err := c.client.GetRankings(context.Background(), req)
	for _, result := range reply.GetResults() {
		r := []*RankingsResult{}
		for _, m := range result.GetRankings() {
			r = append(r, &RankingsResult{Value: m.GetValue(), Count: m.GetCount()})
		}
		ret = append(ret, r)
	}
	return ret, nil
}

// GetCardinality queries the sketch for the cardinality of items.
func (c *Client) GetCardinality(name string) (int64, error) {
	rs := pb.Sketch{Name: &name, Type: &typeCard}
	reply, err := c.client.GetCardinality(context.Background(), &pb.GetRequest{Sketches: []*pb.Sketch{&rs}})
	if err != nil {
		return 0, err
	}
	return reply.GetResults()[0].GetCardinality(), nil
}

// GetMultiCardinality queries multiple sketches for the cardinality of items.
func (c *Client) GetMultiCardinality(names []string) (ret []int64, err error) {
	req := &pb.GetRequest{}
	for i := range names {
		req.Sketches = append(req.Sketches, &pb.Sketch{Name: &names[i], Type: &typeCard})
	}
	reply, err := c.client.GetCardinality(context.Background(), req)

	if err != nil {
		return nil, err
	}
	for _, result := range reply.GetResults() {
		ret = append(ret, result.GetCardinality())
	}
	return ret, nil
}
