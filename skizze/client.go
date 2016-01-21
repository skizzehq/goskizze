package skizze

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/skizzehq/goskizze/datamodel"
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
func (c *Client) ListDomains() (ret []*Domain, err error) {
	reply, err := c.client.ListDomains(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, err
	}
	for _, name := range reply.GetName() {
		ret = append(ret, &Domain{Name: name})
	}
	return ret, err
}
