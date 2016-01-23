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

	lastRequest interface{}
	nextReply   interface{}
	nextError   error
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

func (f *fakeSkizze) List(ctx context.Context, in *pb.ListRequest) (*pb.ListReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.ListReply), f.nextError
}

func (f *fakeSkizze) ListAll(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	return f.nextReply.(*pb.ListReply), f.nextError
}

func (f *fakeSkizze) ListDomains(ctx context.Context, in *pb.Empty) (*pb.ListDomainsReply, error) {
	return f.nextReply.(*pb.ListDomainsReply), f.nextError
}

func (f *fakeSkizze) CreateDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Domain), f.nextError
}

func (f *fakeSkizze) DeleteDomain(ctx context.Context, in *pb.Domain) (*pb.Empty, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Empty), f.nextError
}

func (f *fakeSkizze) GetDomain(ctx context.Context, in *pb.Domain) (*pb.Domain, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Domain), f.nextError
}

func (f *fakeSkizze) CreateSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Sketch), f.nextError
}

func (f *fakeSkizze) DeleteSketch(ctx context.Context, in *pb.Sketch) (*pb.Empty, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Empty), f.nextError
}

func (f *fakeSkizze) GetSketch(ctx context.Context, in *pb.Sketch) (*pb.Sketch, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.Sketch), f.nextError

}

func (f *fakeSkizze) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.AddReply), f.nextError
}

func (f *fakeSkizze) GetMembership(ctx context.Context, in *pb.GetRequest) (*pb.GetMembershipReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.GetMembershipReply), f.nextError
}

func (f *fakeSkizze) GetFrequency(ctx context.Context, in *pb.GetRequest) (*pb.GetFrequencyReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.GetFrequencyReply), f.nextError
}

func (f *fakeSkizze) GetCardinality(ctx context.Context, in *pb.GetRequest) (*pb.GetCardinalityReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.GetCardinalityReply), f.nextError
}

func (f *fakeSkizze) GetRank(ctx context.Context, in *pb.GetRequest) (*pb.GetRankReply, error) {
	f.lastRequest = in
	return f.nextReply.(*pb.GetRankReply), f.nextError
}
