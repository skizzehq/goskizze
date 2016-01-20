package skizze_test

import (
	"net"
	"strconv"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/skizzehq/goskizze/datamodel"
)

type fakeSkizze struct {
	address string
	ready   <-chan bool
	server  *grpc.Server

	nextReply interface{}
	nextError error
}

var port int32 = 6100

func newFakeSkizze() *fakeSkizze {
	p := atomic.AddInt32(&port, 1)

	ready := make(chan bool, 1)
	address := ":" + strconv.Itoa(int(p))
	fs := &fakeSkizze{
		address: address,
		ready:   ready,
	}

	go func() {
		tries := 0
		for {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				if tries < 5 {
					tries++
					continue
				}
				panic(err)
			}

			fs.server = grpc.NewServer()
			pb.RegisterSkizzeServer(fs.server, fs)
			ready <- true
			fs.server.Serve(listener)
			listener.Close()
		}
	}()

	return fs
}

func (f *fakeSkizze) CreateSnapshot(ctx context.Context, in *pb.CreateSnapshotRequest) (*pb.CreateSnapshotReply, error) {
	return f.nextReply.(*pb.CreateSnapshotReply), f.nextError
}

func (f *fakeSkizze) GetSnapshot(ctx context.Context, in *pb.GetSnapshotRequest) (*pb.GetSnapshotReply, error) {
	return f.nextReply.(*pb.GetSnapshotReply), f.nextError
}

func (*fakeSkizze) List(ctx context.Context, in *pb.ListRequest) (*pb.ListReply, error) {
	return nil, nil
}

func (*fakeSkizze) ListAll(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	return nil, nil
}

func (*fakeSkizze) ListDomains(ctx context.Context, in *pb.Empty) (*pb.ListDomainsReply, error) {
	return nil, nil
}

func (*fakeSkizze) SetDefaults(ctx context.Context, in *pb.Defaults) (*pb.Defaults, error) {
	return nil, nil
}

func (*fakeSkizze) GetDefaults(ctx context.Context, in *pb.Empty) (*pb.Defaults, error) {
	return nil, nil
}

func (*fakeSkizze) CreateDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	return nil, nil
}

func (*fakeSkizze) DeleteDomain(ctx context.Context, in *pb.Domain) (*pb.Empty, error) {
	return nil, nil
}

func (*fakeSkizze) GetDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	return nil, nil
}

func (*fakeSkizze) CreateSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	return nil, nil
}

func (*fakeSkizze) DeleteSketch(ctx context.Context, in *pb.Sketch) (*pb.Empty, error) {
	return nil, nil
}

func (*fakeSkizze) GetSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	return nil, nil
}

func (*fakeSkizze) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	return nil, nil
}

func (*fakeSkizze) GetMembership(ctx context.Context, in *pb.GetRequest) (*pb.GetMembershipReply, error) {
	return nil, nil
}

func (*fakeSkizze) GetFrequency(ctx context.Context, in *pb.GetRequest) (*pb.GetFrequencyReply, error) {
	return nil, nil
}

func (*fakeSkizze) GetCardinality(ctx context.Context, in *pb.GetRequest) (*pb.GetCardinalityReply, error) {
	return nil, nil
}

func (*fakeSkizze) GetRank(ctx context.Context, in *pb.GetRequest) (*pb.GetRankReply, error) {
	return nil, nil
}
