package framework

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	metamgr "github.com/basho/bletchley/metadata_manager"
	"github.com/golang/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/basho/bletchley/framework/riak_node_states"

//"github.com/basho/bletchley/common"
)

const (
	OFFER_INTERVAL float64 = 5
)

type SchedulerCore struct {
	lock                *sync.Mutex
	frameworkName       string
	clusters            map[string]*FrameworkRiakCluster
	schedulerHTTPServer *SchedulerHTTPServer
	mgr             *metamgr.MetadataManager
	schedulerIpAddr string
	frnDict				map[string]*FrameworkRiakNode
}

func NewSchedulerCore(schedulerHostname string, frameworkName string, mgr *metamgr.MetadataManager, schedulerIpAddr string) *SchedulerCore {
	scheduler := &SchedulerCore{
		lock:            &sync.Mutex{},
		frameworkName:   frameworkName,
		schedulerIpAddr: schedulerIpAddr,
		clusters:        make(map[string]*FrameworkRiakCluster),
		mgr: mgr,
		frnDict: make(map[string]*FrameworkRiakNode),
	}
	scheduler.schedulerHTTPServer = ServeExecutorArtifact(scheduler, schedulerHostname)
	return scheduler
}

/*

func (sched *SchedulerCore) handleStatusUpdate(msg statusUpdateCast) {
	targetTask, assigned := sched.targetTasksSubs[msg.status.TaskId.GetValue()]
	if assigned {
		targetTask.UpdateStatus(msg.status)
	}
	// We should probably garbage collect the internal task state dictionary
	// But, for now just collect them all -- memory is cheap!
	sched.internalTaskStates[msg.status.TaskId.GetValue()] = msg.status
}
*/

/*
func (sched *SchedulerCore) handleResourceOffers(mesosOffers []*mesos.Offer) {
log.Debugf("Received resource offers: %v", mesosOffers)
launchPlan := make(map[string][]scheduleTask)
outstandingOffers := make(map[string]*mesos.Offer)
for _, offer := range mesosOffers {
	outstandingOffers[offer.Id.GetValue()] = offer
	launchPlan[offer.Id.GetValue()] = []scheduleTask{}
}

for {
	select {
	case request := <-sched.resourceOffersRescinded:
		{
			delete(outstandingOffers, request.offerId.GetValue())
			for _, scheduledTask := range launchPlan[request.offerId.GetValue()] {
				scheduledTask.replyChannel <- false
			}
			delete(launchPlan, request.offerId.GetValue())
		}
	case request := <-sched.outstandingTasks:
		{
			// This actually works, surprisingly enough
			// In order to add multi-task constraints, we need to know what tasks are related to one another
			// and then bucket them appropriately

			// Right now, it fills up individual hosts
			// This is "good enough" (IMHO) for the  demo
			log.Infof("Got asked to schedule outstanding task: %v\n", request)
			for key, offer := range outstandingOffers {
				tmpResources := offer.Resources
				var resourceAsk *mesos.Resource
				var success bool
				asks := []*mesos.Resource{}
				for _, filter := range request.Filters {
					tmpResources, resourceAsk, success = filter(tmpResources)
					if !success {
						break
					}
					asks = append(asks, resourceAsk)
				}
				if success {
					// The new reduced version of the resources
					outstandingOffers[key].Resources = tmpResources
					request.TaskInfo.SlaveId = outstandingOffers[key].SlaveId
					request.TaskInfo.Resources = asks
					launchPlan[key] = append(launchPlan[key], request)
				} else {
					request.replyChannel <- false
				}

			}

		}
	default:
		{
			offerIDs := []*mesos.OfferID{}
			tasks := []*mesos.TaskInfo{}
			for offerID, launchPlanTasks := range launchPlan {
				offerIDs = append(offerIDs, outstandingOffers[offerID].Id)
				for _, task := range launchPlanTasks {
					sched.subscriptionLock.Lock()
					sched.targetTasksSubs[task.TaskInfo.TaskId.GetValue()] = task.TargetTask
					sched.subscriptionLock.Unlock()
					task.replyChannel <- true
					tasks = append(tasks, task.TaskInfo)
				}
			}
			log.Infof("Launching %v task(s) using offerID(s): %v\n", len(tasks), offerIDs)
			sched.driver.LaunchTasks(offerIDs, tasks, &mesos.Filters{RefuseSeconds: proto.Float64(OFFER_INTERVAL)})
			log.Info("No outstanding tasks")
			return
		}
	}
}
}
*/
/*
func (sched *SchedulerCore) SchedulingLoop() {
for {
	select {
	case offers := <-sched.resourceOffers:
		sched.handleResourceOffers(offers.offers)
	// This number is chosen
	case <-time.After(time.Duration(3*OFFER_INTERVAL) * time.Second):
		{
			log.Info("No resource offers received")
			select {
			case request := <-sched.outstandingTasks:
				{
					log.Info("received outstanding tasks during no offer period: ", request)
					request.replyChannel <- false
				}
			default:
				log.Info("Received no outstanding tasks during no offer period")
			}
		}

	}
}
*/
/*
func (sched *SchedulerCore) handleSubChange(subChange taskStateSubscribe) {
	log.Info("Changing subscription: ", subChange)
	sched.subscriptionLock.Lock()
	defer sched.subscriptionLock.Unlock()
	switch subChange.subscriptionChangeType {
	case subscribe:
		{
			// This should trigger a reconcilation
			_, assigned := sched.targetTasksSubs[subChange.taskID]
			if assigned {
				panic("Only one task be assigned to a task ID at a time")
			} else {
				sched.targetTasksSubs[subChange.taskID] = subChange.targetTask
				sched.TriggerReconcilation(subChange.taskID)
			}
		}
	case unsubscribe:
		{
			delete(sched.targetTasksSubs, subChange.taskID)
		}
	}
}
*/
/*
func (sched *SchedulerCore) MesosLoop() {
	initialRegistration := <-sched.registered
	log.Info("Scheduler routine registered: ", initialRegistration.frameworkId, initialRegistration.masterInfo)

	go sched.SchedulingLoop()
	for {
		select {
		case msg := <-sched.reregistered:
			{
				log.Info("Scheduler routine reregistered: ", msg.masterInfo)
			}
		case msg := <-sched.statusUpdate:
			{
				sched.handleStatusUpdate(msg)
			}
		case subChange := <-sched.taskStateSubscribe:
			{
				sched.handleSubChange(subChange)

			}
		}
	}
}
*/

/*
func (sched *SchedulerCore) TriggerReconcilation(taskID string) {
	ts := mesos.TaskState_TASK_ERROR
	task := &mesos.TaskStatus{
		TaskId:  &mesos.TaskID{Value: proto.String(taskID)},
		State:   &ts,
		SlaveId: &mesos.SlaveID{Value: proto.String("")}, // Slave ID isn't required
	}
	taskStatuses := []*mesos.TaskStatus{task}
	sched.driver.ReconcileTasks(taskStatuses)
}
*/

// This is an add cluster callback from the metadata manager
func (sc *SchedulerCore) AddCluster(zkNode *metamgr.ZkNode) metamgr.MetadataManagerCluster {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	frc := NewFrameworkRiakCluster()
	frc.sc = sc
	frc.zkNode = zkNode
	err := json.Unmarshal(zkNode.GetData(), &frc)
	if err != nil {
		log.Panic("Error getting node: ", err)
	}
	sc.clusters[frc.Name] = frc
	return frc
}
func (sc *SchedulerCore) GetCluster(name string) metamgr.MetadataManagerCluster {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	return sc.clusters[name]
}

// Should basically just be a callback - DO NOT change state
func (sc SchedulerCore) NewCluster(zkNode *metamgr.ZkNode, name string) metamgr.MetadataManagerCluster {
	frc := &FrameworkRiakCluster{
		zkNode: zkNode,
		nodes:  make(map[string]*FrameworkRiakNode),
		Name:   name,
	}
	return frc
}

func (sc *SchedulerCore) setupMetadataManager() {
	sc.mgr.SetupFramework(sc.schedulerHTTPServer.URI, sc)
	for _, cluster := range sc.clusters {
		for _, riakNode := range cluster.nodes {
			if riakNode.TaskStatus != nil {
				log.Info("Populate FRN dict: %v -> %v", riakNode.TaskStatus.TaskId.GetValue(), riakNode)
				sc.frnDict[riakNode.TaskStatus.TaskId.GetValue()] = riakNode
			}
		}
	}
}
func (sc *SchedulerCore) Run(mesosMaster string) {
	frameworkId := &mesos.FrameworkID{
		Value: proto.String(sc.frameworkName),
	}
	// TODO: Get "Real" credentials here
	cred := (*mesos.Credential)(nil)
	bindingAddress := parseIP(sc.schedulerIpAddr)
	fwinfo := &mesos.FrameworkInfo{
		User:            proto.String("sargun"), // Mesos-go will fill in user.
		Name:            proto.String("Test Framework (Go)"),
		Id:              frameworkId,
		FailoverTimeout: proto.Float64(86400),
	}
	config := sched.DriverConfig{
		Scheduler:      sc,
		Framework:      fwinfo,
		Master:         mesosMaster,
		Credential:     cred,
		BindingAddress: bindingAddress,
		//	WithAuthContext: func(ctx context.Context) context.Context {
		//		ctx = auth.WithLoginProvider(ctx, *authProvider)
		//		ctx = sasl.WithBindingAddress(ctx, bindingAddress)
		//		return ctx
		//	},
	}
	driver, err := sched.NewMesosSchedulerDriver(config)
	if err != nil {
		log.Error("Unable to create a SchedulerDriver ", err.Error())
	}
	sc.setupMetadataManager()
	if stat, err := driver.Run(); err != nil {
		log.Infof("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
	}
}
func (sc *SchedulerCore) reconcileTasks(driver sched.SchedulerDriver) {
	tasksToReconcile := []*mesos.TaskStatus{}
	for _, cluster := range sc.clusters {
		for _, riakNode := range cluster.nodes {
			tasksToReconcile = append(tasksToReconcile, riakNode.GetTaskStatus())
		}
	}
	driver.ReconcileTasks(tasksToReconcile)
}
func (sc *SchedulerCore) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Framework registered")
	sc.reconcileTasks(driver)
}

func (sc *SchedulerCore) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	//go NewTargetTask(*sched).Loop()
	// We don't actually handle this correctly
	log.Error("Framework reregistered")
	sc.reconcileTasks(driver)
}
func (sc *SchedulerCore) Disconnected(sched.SchedulerDriver) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Error("Framework disconnected")
}

func (sc *SchedulerCore) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Received resource offers")
	executorUris := []*mesos.CommandInfo_URI{}
	executorUris = append(executorUris,
		&mesos.CommandInfo_URI{Value: &(sc.schedulerHTTPServer.hostURI), Executable: proto.Bool(true)})
	launchTasks := []*mesos.TaskInfo{}
	for _, cluster := range sc.clusters {
		for _, riakNode := range cluster.nodes {
			if riakNode.NeedsToBeScheduled() {
				log.Infof("Adding Riak node for scheduling: %+v", riakNode)
				// We need to schedule this task I guess?
				riakNode.generation = riakNode.generation + 1
				task_id := riakNode.NewID()
				exec := &mesos.ExecutorInfo{
					//No idea is this is the "right" way to do it, but I think so?
					ExecutorId: util.NewExecutorID(task_id),
					Name:       proto.String("Test Executor (Go)"),
					Source:     proto.String("Riak Mesos Framework (Go)"),
					Command: &mesos.CommandInfo{
						Value:     proto.String(sc.schedulerHTTPServer.executorName),
						Uris:      executorUris,
						Shell:     proto.Bool(false),
						Arguments: []string{sc.schedulerHTTPServer.executorName, "-taskid", task_id},
					},
				}

				taskId := &mesos.TaskID{
					Value: proto.String(task_id),
				}
				taskInfo := &mesos.TaskInfo{
					Name:     proto.String(task_id),
					TaskId:   taskId,
					SlaveId:  offers[0].SlaveId,
					Executor: exec,
					Resources: []*mesos.Resource{
						util.NewScalarResource("mem", 1),
					},
					Data: []byte{'h', 'e', 'l', 'l', 'o'},
				}
				ts := mesos.TaskState_TASK_ERROR
				taskStatus := &mesos.TaskStatus{
					TaskId:  &mesos.TaskID{Value: proto.String(task_id)},
					State:   &ts,
					SlaveId: &mesos.SlaveID{Value: proto.String("")}, // Slave ID isn't required
				}
				sc.frnDict[task_id] = riakNode
				riakNode.CurrentState = riak_node_states.Starting
				riakNode.TaskStatus = taskStatus
				riakNode.Persist()
				launchTasks = append(launchTasks, taskInfo)
			}
		}
	}
	offerIDs := make([]*mesos.OfferID, len(offers))
	for idx, offer := range offers {
		offerIDs[idx] = offer.Id
	}
	log.Info("Launching Tasks: ", launchTasks)
	driver.LaunchTasks(offerIDs, launchTasks, &mesos.Filters{RefuseSeconds: proto.Float64(OFFER_INTERVAL)})
}
func (sc *SchedulerCore) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	riak_node, assigned := sc.frnDict[status.TaskId.GetValue()]
	if assigned {
		log.Info("Received status updates: ", status)
		log.Info("Riak Node: ", riak_node)
		riak_node.handleStatusUpdate(status)
		riak_node.Persist()
	} else {
		log.Error("Received status update for unknown job: ", status)
	}

}

func (sc *SchedulerCore) OfferRescinded(driver sched.SchedulerDriver, offerID *mesos.OfferID) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Offer rescinded from Mesos")
}

func (sc *SchedulerCore) FrameworkMessage(driver sched.SchedulerDriver, executorID *mesos.ExecutorID, slaveID *mesos.SlaveID, message string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Got unknown framework message %v")
}

// TODO: Write handler
func (sc *SchedulerCore) SlaveLost(sched.SchedulerDriver, *mesos.SlaveID) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Slave Lost")
}

// TODO: Write handler
func (sc *SchedulerCore) ExecutorLost(sched.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, int) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Executor Lost")
}

func (sc *SchedulerCore) Error(driver sched.SchedulerDriver, err string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Info("Scheduler received error:", err)
}

// Old:

/*

// Private "internal" structs
type registeredCast struct {
	frameworkId *mesos.FrameworkID
	masterInfo  *mesos.MasterInfo
}
type reregisteredCast struct {
	masterInfo *mesos.MasterInfo
}
type statusUpdateCast struct {
	status *mesos.TaskStatus
}

type SubscriptionChangeType int

const (
	subscribe   SubscriptionChangeType = iota
	unsubscribe                        = iota
)

type taskStateSubscribe struct {
	targetTask             *TargetTask
	taskID                 string
	subscriptionChangeType SubscriptionChangeType
}

type resourceOffers struct {
	offers []*mesos.Offer
}

type resourceOffersRescinded struct {
	offerId *mesos.OfferID
}

type scheduleTask struct {
	TaskInfo     *mesos.TaskInfo
	TargetTask   *TargetTask
	Filters      []common.ResourceAsker
	replyChannel chan bool
}
*/
