package mq

import (
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

// Test the msgbase message creation using generics
func TestServerConnection(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Create a new nats client from the environment variables
	natsClient, natsClientError := new(mq.Nats).Init()

	if natsClientError != nil {
		t.Skip("nats not configured for test")
	}

	// Connect to the nats server
	err := natsClient.Connect()

	if err != nil {
		zap.L().Error("failed to connect to nats server")
		t.Fatal()
	}

}
