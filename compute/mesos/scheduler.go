package mesos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/api/v0/auth"
	"github.com/mesos/mesos-go/api/v0/mesosproto"
	"github.com/mesos/mesos-go/api/v0/scheduler"
	//"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
)

// Settings holds configuration values for the scheduler
type Settings struct {
	Master           string
	FrameworkID      string
	CredentialFile   string
	Name             string
	User             string
	MessengerAddress string
	MessengerPort    uint16
	Checkpoint       bool
	FailoverTimeout  float64
}

// Scheduler holds the structure of Funnel's Mesos Scheduler
type Scheduler struct {
	conf  *Settings
	log   *logger.Logger
	driver mesossched.SchedulerDriver	
}

// NewScheduler returns a new instance of Funnel's Mesos scheduler.
func NewScheduler(conf *Settings, log *logger.Logger) (*Scheduler, error) {
	scheduler := &Scheduler{
		conf: conf,
		log:  log,
	}

	publishedAddr := net.ParseIP(conf.MessengerAddress)
	bindingPort := conf.MessengerPort
	credential :=  &mesosproto.Credential{}
	var principal *string
	if credential != nil {
		principal = credential.Principal
	}
	
	getAuthContext := func(ctx context.Context) context.Context {
		return auth.WithLoginProvider(ctx, "SASL")
	}

	driver := mesossched.NewMesosSchedulerDriver(mesossched.DriverConfig{
		Master: settings.Master,
		Framework: &mesosproto.FrameworkInfo{
			Id:              &mesosproto.FrameworkID{
				Value: conf.FrameworkID,
			},
			Name:            conf.Name,
			User:            conf.User,
			Checkpoint:      proto.Bool(conf.Checkpoint),
			FailoverTimeout: proto.Float64(conf.FailoverTimeout),
			Principal:       principal,
		},
		Scheduler:        scheduler,
		BindingAddress:   net.ParseIP("0.0.0.0"),
		PublishedAddress: publishedAddr,
		BindingPort:      bindingPort,
		Credential:       credential,
		WithAuthContext:  getAuthContext,
	})

	scheduler.driver = driver
	return scheduler, nil
}

func (s *Scheduler) Run(pctx context.Context) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()
	s.log.Info("Starting funnel mesos scheduler")

	go func() {
		<-ctx.Done()
		driver.Stop(false)
	}()

	if status, err := driver.Run(); err != nil {
		s.log.Error("Mesos framework stopped", "error", err, "status", status.String())
	}

	return
}

// WriteEvent writes an event to the compute backend.
// Currently, handles TASK_CREATED and TASK_STATE events. 
func (s *Scheduler) WriteEvent(ctx context.Context, ev *events.Event) error {
	switch ev.Type {
	case events.Type_TASK_CREATED:
		return s.Submit(ev.GetTask())

	case events.Type_TASK_STATE:
		if ev.GetState() == tes.State_CANCELED {
			return s.Cancel(ctx, ev.Id)
		}
	}
	return nil
}

// Submit submits a task to the mesos scheduler.
func (s *Scheduler) Submit(task *tes.Task) error {
	return nil
}

// Cancel sends a cancel signal to a task managed by mesos.
func (s *Scheduler) Cancel(taskID string) error {
	_, err = s.driver.KillTask(&mesosproto.TaskID{
		Value: taskID,
	})
	return err
}

// Registered is called when the Scheduler is Registered
func (s *Scheduler) Registered(driver mesossched.SchedulerDriver, frameworkID *mesosproto.FrameworkID, masterInfo *mesosproto.MasterInfo) {
	s.log.Debug("Framework registered with master", "framework_id", frameworkID.GetValue(), "master_id", masterInfo.GetId(), "master", masterInfo.GetHostname())
}

// Reregistered is called when the Scheduler is Reregistered
func (s *Scheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	s.log.Debug("Framework re-registered with master", "master_id", masterInfo.GetId(), "master", masterInfo.GetHostname())
}

// Disconnected is called when the Scheduler is Disconnected
func (s *Scheduler) Disconnected(sched.SchedulerDriver) {
	s.log.Debug("Framework disconnected with master")
}

// ResourceOffers handles the Resource Offers
func (sched *ExampleScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	s.log.Debug("Recieved offers", "offers", len(offers))
}

// StatusUpdate takes care of updating the status
func (sched *ExampleScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	s.log.Debug("Recieved status update", "task_id", status.TaskId.GetValue(), "status", status.State.String())
}

// OfferRescinded is invoked when an offer is no longer valid.
func (sched *ExampleScheduler) OfferRescinded(_ sched.SchedulerDriver, oid *mesos.OfferID) {
	s.log.Debug("Offer rescinded", "offer_id", oid)
}

// FrameworkMessage is invoked when an executor sends a message.
func (sched *ExampleScheduler) FrameworkMessage(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	s.log.Debug("Recieved framework message", "msg", msg)
}

// SlaveLost is invoked when a slave has been determined unreachable.
func (sched *ExampleScheduler) SlaveLost(_ sched.SchedulerDriver, sid *mesos.SlaveID) {
	s.log.Debug("Slave lost", "slave_id", sid)
}

// ExecutorLost is invoked when an executor has exited/terminated.
func (sched *ExampleScheduler) ExecutorLost(_ sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	s.log.Debug("Executor lost", "slave_id", sid, "executor_id", eid)
}

// Error is invoked when there is an unrecoverable error in the scheduler or scheduler driver.
func (sched *ExampleScheduler) Error(_ sched.SchedulerDriver, err string) {
	s.log.Debug("Recieved an error", "error", err)
}
