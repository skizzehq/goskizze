package skizze

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/skizzehq/goskizze/datamodel"
)

type fakeSkizze struct {
	address string
	ready   <-chan bool
	server  *grpc.Server
}

func newFakeSkizze() *fakeSkizze {
	ready := make(chan bool, 1)
	address := ":6111"
	fs := &fakeSkizze{
		address: address,
		ready:   ready,
	}

	go func() {
		for {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				panic(err)
			}
			defer listener.Close()

			fs.server = grpc.NewServer()
			pb.RegisterSkizzeServer(fs.server, &fakeSkizze{})
			ready <- true
			fs.server.Serve(listener)
		}
	}()

	return fs
}

func (*fakeSkizze) CreateSnapshot(ctx context.Context, in *pb.CreateSnapshotRequest) (*pb.CreateSnapshotReply, error) {
	return nil, nil
}

func (*fakeSkizze) GetSnapshot(ctx context.Context, in *pb.GetSnapshotRequest) (*pb.GetSnapshotReply, error) {
	return nil, nil
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
