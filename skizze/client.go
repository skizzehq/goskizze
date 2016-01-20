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
