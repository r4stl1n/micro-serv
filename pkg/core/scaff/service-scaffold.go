package scaff

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"go.uber.org/zap"
	"time"
)

type ServiceScaffold struct {
	nats *mq.Nats

	scaffoldContext *ScaffoldContext
	serviceMan      micro.Service
	serviceGroup    micro.Group
	serviceConfig   micro.Config

	taskScheduler *gocron.Scheduler
}

// Init create a service scaffold
func (s *ServiceScaffold) Init(serviceName string, version string, description string) (*ServiceScaffold, error) {

	// Create the nats connection
	nats, natsError := new(mq.Nats).Init()

	if natsError != nil {
		return nil, natsError
	}

	// Create the microservice configuration
	config := micro.Config{
		Name:        serviceName,
		Version:     version,
		Description: description,
		// DoneHandler can be set to customize behavior on stopping a service.
		DoneHandler: func(srv micro.Service) {
			info := srv.Info()
			zap.L().Info("stopped service", zap.String("name", info.Name), zap.String("id", info.ID))
		},
		// ErrorHandler can be used to customize behavior on service execution error.
		ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
			info := srv.Info()
			zap.L().Error("service encountered error", zap.String("name", info.Name),
				zap.String("subject", err.Subject), zap.String("description", err.Description))
		},
	}

	*s = ServiceScaffold{
		nats:          nats,
		serviceConfig: config,
		taskScheduler: gocron.NewScheduler(time.UTC),
	}

	return s, nil
}

func (s *ServiceScaffold) IsRunning() bool {
	return !s.serviceMan.Stopped()
}

// AddHandler adds the new handler to the scaff
func (s *ServiceScaffold) AddHandler(context any, subject string,
	handler func(context any, scaffoldContext *ScaffoldContext, request micro.Request)) error {

	if s.serviceMan.Stopped() {
		return fmt.Errorf("service needs to be started first before adding handlers")
	}

	var endpointError error

	endpointError = s.serviceGroup.AddEndpoint(subject,
		ScaffoldContextHelper(context, s.scaffoldContext, handler))

	return endpointError
}

// AddRepeatTask adds a task that is routinely called at specified duration
// interval is a specific time string ex. 30s 1m 5m 1h
func (s *ServiceScaffold) AddRepeatTask(context any, interval string,
	taskHandler func(context any, scaffoldContext *ScaffoldContext)) error {

	if s.serviceMan.Stopped() {
		return fmt.Errorf("service needs to be started first before adding tasks")
	}

	t := time.Now()

	if context == nil {
		context = ""
	}

	_, scheduleError := s.taskScheduler.Every(interval).StartAt(
		time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)).SingletonMode().Do(
		taskHandler, context, s.scaffoldContext)

	if scheduleError != nil {
		return scheduleError
	}

	return nil
}

// Start the service with a context,
// Note: context can be any structure or data for use by the handlers
func (s *ServiceScaffold) Start() error {

	// Connect to the nats server
	connectError := s.nats.Connect()

	if connectError != nil {
		return fmt.Errorf("failed to connect to nats service: %s", s.nats.GetConfig().Host)
	}

	zap.L().Info("connected to server", zap.String("server", s.nats.GetConfig().Host))

	// Create the microservice (it starts automatically)
	service, serviceManError := micro.AddService(s.nats.GetConn(), s.serviceConfig)

	if serviceManError != nil {
		return serviceManError
	}

	// Add the service
	s.serviceMan = service

	// Add the service group
	s.serviceGroup = s.serviceMan.AddGroup("v1")

	// Start the task scheduler
	s.taskScheduler.StartAsync()

	// Create the scaffold context
	s.scaffoldContext = new(ScaffoldContext).Init(s.nats, s.serviceMan.Info())

	return nil
}

func (s *ServiceScaffold) Stop() error {

	if s.serviceMan.Stopped() {
		return fmt.Errorf("failed to stop: service is not running")
	}

	s.taskScheduler.StopBlockingChan()

	return s.serviceMan.Stop()
}

func (s *ServiceScaffold) ServiceID() string {
	return s.serviceMan.Info().ID
}

func (s *ServiceScaffold) SendRequest(subject string, data any, result any) error {

	// Send the add symbol request
	return s.nats.SendRequest(subject, data, &result)
}

func (s *ServiceScaffold) SendMessage(subject string, data any) error {

	// Send the add symbol request
	return s.nats.SendMessage(subject, data)
}
