package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	networkingV1alpha3 "github.com/ruiwang47/k8s-istio-client/pkg/apis/networking/v1alpha3"
	clientset "github.com/ruiwang47/k8s-istio-client/pkg/client/clientset/versioned"
	networkingScheme "github.com/ruiwang47/k8s-istio-client/pkg/client/clientset/versioned/scheme"
	informers "github.com/ruiwang47/k8s-istio-client/pkg/client/informers/externalversions/networking/v1alpha3"
	listers "github.com/ruiwang47/k8s-istio-client/pkg/client/listers/networking/v1alpha3"
)

const controllerAgentName = "virtualservices-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a virtualservices is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a virtualservices
	// is synced successfully
	MessageResourceSynced = "virtualservices synced successfully"
)

// Controller is the controller implementation for virtualservices resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// virtualservicesClientset is a clientset for our own API group
	virtualservicesClientset clientset.Interface

	virtualservicesLister listers.VirtualServiceLister
	virtualservicesSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new virtualservices controller
func NewController(
	kubeclientset kubernetes.Interface,
	virtualservicesClientset clientset.Interface,
	virtualservicesInformer informers.VirtualServiceInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for virtualservices controller types.
	utilruntime.Must(networkingScheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:            kubeclientset,
		virtualservicesClientset: virtualservicesClientset,
		virtualservicesLister:    virtualservicesInformer.Lister(),
		virtualservicesSynced:    virtualservicesInformer.Informer().HasSynced,
		workqueue:                workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks"),
		recorder:                 recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when virtualservices resources change
	virtualservicesInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueueVirtualServices,
		UpdateFunc: controller.enqueueVirtualServicesForUpdate,
		DeleteFunc: controller.enqueueVirtualServicesForDelete,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting virtualservices control loop")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.virtualservicesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch workers to process virtualservices resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// virtualservices resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the virtualservices resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the virtualservices resource with this namespace/name
	virtualService, err := c.virtualservicesLister.VirtualServices(namespace).Get(name)
	if err != nil {
		// The virtualservices resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			glog.Warningf("virtualservices: %s/%s does not exist in local cache, will delete it from Neutron ...",
				namespace, name)

			glog.Infof("[Neutron] Deleting virtualService: %s/%s ...", namespace, name)

			return nil
		}

		runtime.HandleError(fmt.Errorf("failed to list virtualService by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("[Neutron] Try to process virtualService: %#v ...", virtualService)

	c.recorder.Event(virtualService, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) enqueueVirtualServices(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	glog.Info("add ", key)
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueVirtualServicesForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	glog.Info("delete ", key)
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueVirtualServicesForUpdate(old, new interface{}) {
	oldVs := old.(*networkingV1alpha3.VirtualService)
	newVs := new.(*networkingV1alpha3.VirtualService)
	if oldVs.ResourceVersion == newVs.ResourceVersion {
		// Periodic resync will send update events for all known virtualservices.
		// Two different versions of the same virtualservices will always have different RVs.
		return
	}
	glog.Info("update ", old, new)
	c.enqueueVirtualServices(new)
}
