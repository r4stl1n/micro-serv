package mq

import (
	"fmt"
	"os"
)

// NatsConfig configuration structure for the nats mq
type NatsConfig struct {
	Host string
	User string
	Pass string
}

// FromEnv creates a nats config using environmental variables
func (n *NatsConfig) FromEnv() (*NatsConfig, error) {

	natsHost, natsHostOk := os.LookupEnv("NATS_HOST")
	natsUser, natsUserOk := os.LookupEnv("NATS_USER")
	natsPass, natsPassOk := os.LookupEnv("NATS_PASS")

	if !natsHostOk {
		return nil, fmt.Errorf("NATS_HOST enviornment variable not set")
	}

	if !natsUserOk {
		return nil, fmt.Errorf("NATS_USER enviornment variable not set")
	}

	if !natsPassOk {
		return nil, fmt.Errorf("NATS_PASS enviornment variable not set")
	}

	*n = NatsConfig{
		Host: natsHost,
		User: natsUser,
		Pass: natsPass,
	}

	return n, nil
}

// IsSet checks if the configuration has been set
// Note: All fields required
func (n *NatsConfig) IsSet() bool {
	return n.Host != "" && n.User != "" && n.Pass != ""
}
