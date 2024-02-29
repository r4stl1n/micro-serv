package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"reflect"
	"time"
)

// Nats creates and manages the nats message queue connection
type Nats struct {
	config *NatsConfig
	conn   *nats.Conn
	error  error
}

// Init create a new nats manager with a given config
func (n *Nats) Init() (*Nats, error) {

	config, configError := new(NatsConfig).FromEnv()

	if configError != nil {
		return nil, configError
	}

	*n = Nats{
		config: config,
	}

	return n, nil
}

// GetConfig get the nats config
func (n *Nats) GetConfig() NatsConfig {
	return *n.config
}

// Connect to the nats server
func (n *Nats) Connect() error {

	// Check if our nats configuration is set properly
	if !n.config.IsSet() {
		return fmt.Errorf("nats not configured")
	}

	// Create the nats connection string
	descriptor := fmt.Sprintf("mq://%s:%s@%s", n.config.User, n.config.Pass, n.config.Host)

	// Attempt to connect to nats
	n.conn, n.error = nats.Connect(descriptor, nats.ReconnectWait(5*time.Second), nats.Timeout(5*time.Second))
	if n.error != nil {
		return n.error
	}

	zap.L().Info("connected to nats server", zap.String("host", n.config.Host))

	return n.error
}

// Disconnect from the nats server
func (n *Nats) Disconnect() {
	n.conn.Close()

	zap.L().Info("disconnected from nats server", zap.String("host", n.config.Host))
}

// GetConn retrieves the raw nats connection
func (n *Nats) GetConn() *nats.Conn {
	return n.conn
}

// SendRequest sends a request and captures the error
func (n *Nats) SendRequest(subject string, data any, result any) error {

	// Create the add symbol request
	marshallBytes, marshallError := json.Marshal(&data)

	if marshallError != nil {
		return marshallError
	}

	zap.L().Info("sending nats request", zap.String("subject", subject),
		zap.Int("size", int(reflect.TypeOf(data).Size())))

	// Send the request over our nats connection
	response, responseError := n.conn.Request(subject, marshallBytes, 30*time.Second)

	if responseError != nil {
		return responseError
	}

	// Check if there is a nats service error
	if response.Header.Get("Nats-Service-Error") != "" {
		return fmt.Errorf(response.Header.Get("Nats-Service-Error"))
	}

	zap.L().Info("received nats response", zap.String("subject", subject),
		zap.Int("size", len(response.Data)))

	// Unmarshall the object for the result
	unmarshallError := json.Unmarshal(response.Data, &result)

	if unmarshallError != nil {
		return unmarshallError
	}

	return nil
}

// SendMessage sends a request and captures the error
func (n *Nats) SendMessage(subject string, data any) error {

	// Create the add symbol request
	marshallBytes, marshallError := json.Marshal(&data)

	if marshallError != nil {
		return marshallError
	}

	zap.L().Info("sending nats message", zap.String("subject", subject),
		zap.Int("size", int(reflect.TypeOf(data).Size())))

	return n.conn.Publish(subject, marshallBytes)
}
