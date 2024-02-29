package db

import (
	"fmt"
	"os"
)

// MongoConfig configuration structure for the db db
type MongoConfig struct {
	Host string
	User string
	Pass string
}

// FromEnv creates a nats config using environmental variables
func (n *MongoConfig) FromEnv() (*MongoConfig, error) {

	natsHost, natsHostOk := os.LookupEnv("MONGO_HOST")
	natsUser, natsUserOk := os.LookupEnv("MONGO_USER")
	natsPass, natsPassOk := os.LookupEnv("MONGO_PASS")

	if !natsHostOk {
		return nil, fmt.Errorf("MONGO_HOST enviornment variable not set")
	}

	if !natsUserOk {
		return nil, fmt.Errorf("MONGO_USER enviornment variable not set")
	}

	if !natsPassOk {
		return nil, fmt.Errorf("MONGO_PASS enviornment variable not set")
	}

	*n = MongoConfig{
		Host: natsHost,
		User: natsUser,
		Pass: natsPass,
	}

	return n, nil
}

// IsSet checks if the configuration has been set
// Note: All fields required
func (n *MongoConfig) IsSet() bool {
	return n.Host != "" && n.User != "" && n.Pass != ""
}
