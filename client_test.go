package skizze

import (
	"testing"

	//	"golang.org/x/net/context"
	//	"google.golang.org/grpc"

	//	pb "github.com/skizzehq/goskizze/datamodel"
)

func TestDial(t *testing.T) {
	fs := newFakeSkizze()
	<-fs.ready

	c := Dial(fs.address, &Options{Insecure: true})
}
