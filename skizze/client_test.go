package skizze

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"golang.org/x/net/context"
	//	"google.golang.org/grpc"

	//	pb "github.com/skizzehq/goskizze/datamodel"
)

func getClient(t *testing.T) *Client {
	assert := assert.New(t)

	fs := newFakeSkizze()
	<-fs.ready
	defer fs.server.Stop()

	c, err := Dial(fs.address, Options{Insecure: true})
	assert.Nil(err)

	return c
}

func TestDial(t *testing.T) {
	getClient(t)
}
