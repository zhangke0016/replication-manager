// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Copyright 2017 Signal 18 SARL
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <svaroqui@gmail.com>
// This source code is licensed under the GNU General Public License, version 3.

package receiver

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/signal18/replication-manager/graphite/points"
)

func TestStopUDP(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveUDPAddr("udp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		r, err := New("udp://" + addr.String())
		assert.NoError(err)
		addr = r.(*UDP).Addr().(*net.UDPAddr) // listen same port in next iteration
		r.Stop()
	}
}

func TestStopTCP(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		r, err := New("tcp://" + addr.String())
		assert.NoError(err)
		addr = r.(*TCP).Addr().(*net.TCPAddr) // listen same port in next iteration
		r.Stop()
	}
}

func TestStopPickle(t *testing.T) {
	assert := assert.New(t)

	addr, err := net.ResolveTCPAddr("tcp", ":0")
	assert.NoError(err)

	for i := 0; i < 10; i++ {
		r, err := New("pickle://" + addr.String())
		assert.NoError(err)
		addr = r.(*TCP).Addr().(*net.TCPAddr) // listen same port in next iteration
		r.Stop()
	}
}

func TestStopConnectedTCP(t *testing.T) {
	test := newTCPTestCase(t, false)
	defer test.Finish()

	ch := test.rcvChan
	test.Send("hello.world 42.15 1422698155\n")
	time.Sleep(10 * time.Millisecond)

	select {
	case msg := <-ch:
		test.Eq(msg, points.OnePoint("hello.world", 42.15, 1422698155))
	default:
		t.Fatalf("Message #0 not received")
	}

	test.receiver.Stop()
	test.receiver = nil
	time.Sleep(10 * time.Millisecond)

	test.Send("metric.name -72.11 1422698155\n")
	time.Sleep(10 * time.Millisecond)

	select {
	case <-ch:
		t.Fatalf("Message #0 received")
	default:
	}
}
