package receiver

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lomik/go-carbon/cache"
	"github.com/lomik/go-carbon/points"
	"github.com/stretchr/testify/assert"
)

func TestStopUDP(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveUDPAddr("udp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		listener := NewUDP(cache.New())
		assert.NoError(listener.Listen(addr))
		addr = listener.Addr().(*net.UDPAddr) // listen same port in next iteration
		listener.Stop()
	}
}

func TestStopTCP(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		listener := NewTCP(cache.New())
		assert.NoError(listener.Listen(addr))
		addr = listener.Addr().(*net.TCPAddr) // listen same port in next iteration
		listener.Stop()
	}
}

func TestStopPickle(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		listener := NewPickle(cache.New())
		assert.NoError(listener.Listen(addr))
		addr = listener.Addr().(*net.TCPAddr) // listen same port in next iteration
		listener.Stop()
	}
}

func TestStopConnectedTCP(t *testing.T) {
	test := newTCPTestCase(t, false)
	defer test.Finish()

	metric := "hello.world"
	test.Send(fmt.Sprintf("%s 42.15 1422698155\n", metric))
	test.GetEq(metric, points.OnePoint(metric, 42.15, 1422698155))

	test.receiver.Stop()
	test.receiver = nil
	time.Sleep(10 * time.Millisecond)

	test.Send("metric.name -72.11 1422698155\n")

	_, ok := test.Get("metric.name")
	assert.False(t, ok, "Metric was sent despite stopped receiver")
}
