package service

import (
	"fmt"
	"github.com/r4stl1n/micro-serv/internal/core/dummy-service/handlers"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

type DummyService struct {
	scaffold       *scaff.ServiceScaffold
	endpointPrefix string
}

// Init create a new DummyService
func (d *DummyService) Init(endpointPrefix string) (*DummyService, error) {

	// Create a new scaffold instance
	scaffold, scaffoldingError := new(scaff.ServiceScaffold).Init("dummy-service", "1.0.0", "Dummy Service")

	if scaffoldingError != nil {
		return nil, scaffoldingError
	}

	*d = DummyService{
		scaffold:       scaffold,
		endpointPrefix: endpointPrefix,
	}

	return d, nil
}

// StartUp create the handlers and perform any additional first time needs
func (d *DummyService) addHandlers() error {

	echoHandlerError := d.scaffold.AddHandler(nil,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDummyEndpointEcho), handlers.EchoHandler)

	if echoHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDummyEndpointEcho, echoHandlerError)
	}

	return nil
}

// simpleTask repeating task to validate repeat task work
func (d *DummyService) simpleTask(_ any, _ *scaff.ScaffoldContext) {
	zap.L().Info("simpleTask queue")
}

// addTasks adds a simple task that simply repeats every 10 sec
func (d *DummyService) addTasks() error {
	taskError := d.scaffold.AddRepeatTask(nil, "10s", d.simpleTask)

	if taskError != nil {
		return fmt.Errorf("failed to add simpleTask")
	}

	return nil
}

// Run the service
func (d *DummyService) Run() error {

	startError := d.scaffold.Start()

	if startError != nil {
		return fmt.Errorf("failed to start the service: %s", startError)
	}

	_ = d.addTasks()

	return d.addHandlers()
}

// Stop the service
func (d *DummyService) Stop() error {
	stopError := d.scaffold.Stop()

	if stopError != nil {
		return fmt.Errorf("failed to stop the service: %s", stopError)
	}

	return nil
}

// IsRunning returns if the service is running
func (d *DummyService) IsRunning() bool {
	return d.scaffold.IsRunning()
}
