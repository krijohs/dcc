package controller

import (
	"reflect"
	"sync/atomic"
	"time"

	"github.com/krijohs/dcc/pkg/logger"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	resourceType = "namespace"
)

const (
	Add = iota
	Update
	Delete
)

type Config struct {
	KubeConf     string
	ResyncPeriod time.Duration
}

type Controller struct {
	log      logger.Logger
	config   Config
	client   kubernetes.Interface
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	handler  Handlerer
}

type Handlerer interface {
	HandleItem(action int, ns string) error
}

// New initializes a new controller which is used to watch namespaces.
func New(log logger.Logger, config Config, client kubernetes.Interface, handler Handlerer) *Controller {
	c := Controller{log: log, config: config, client: client, handler: handler}

	c.initController()
	return &c
}

type eventItem struct {
	action int
	key    string
}

func (c *Controller) initController() {
	sharedFactory := informers.NewSharedInformerFactory(c.client, c.config.ResyncPeriod)
	c.informer = sharedFactory.Core().V1().Namespaces().Informer()

	c.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	var ev eventItem
	var err error
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ev.key, err = cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				c.log.Errorf("unable to use cache.MetaNamspaceKeyFunc: %v", err)
				return
			}

			ev.action = Add
			c.queue.Add(ev)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if !reflect.DeepEqual(oldObj, newObj) {
				ev.key, err = cache.MetaNamespaceKeyFunc(newObj)
				if err != nil {
					c.log.Errorf("unable to use cache.MetaNamspaceKeyFunc: %v", err)
					return
				}

				ev.action = Update
				c.queue.Add(ev)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ev.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				c.log.Errorf("unable to use cache.DeletetionHandlingMetaNamespaceKeyFunc: %v", err)
				return
			}

			ev.action = Delete
			c.queue.Add(ev)
		},
	})
}

// Watch uses the informer to watch for namespace state changes in the cluster.
func (c *Controller) Watch(stopCh chan struct{}) {
	var hasSynced int32

	// kubernetes utility to handle API crashes
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	go func() {
		time.Sleep(5 * time.Second)
		if atomic.LoadInt32(&hasSynced) == 0 {
			stopCh <- struct{}{}
		}
	}()

	// ensure cache has synced atleast once
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(errors.New("timed out waiting for first time cache to sync"))
		return
	}
	atomic.AddInt32(&hasSynced, 1)

	go wait.Until(func() { c.processQueue() }, time.Second, stopCh)
	<-stopCh
}

func (c *Controller) processQueue() bool {
	item, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(item)

	ei, ok := item.(eventItem)
	if !ok {
		c.log.Errorf("unable to type assert item: %+v from queue to eventItem", item)
		return true
	}

	err := c.handler.HandleItem(ei.action, ei.key)
	switch {
	case err == nil:
		// If no error, tell the queue to stop tracking history for the key.
		// This will reset things like failure counts for per-item rate limiting.
		c.queue.Forget(item)
		return true

	// requeue 10 times
	case c.queue.NumRequeues(item) < 10:
		c.queue.AddRateLimited(item)

	default:
		c.queue.Forget(item)
		c.log.Errorf("unable to process item: %+v , forgetting about it", item)
		utilruntime.HandleError(errors.Wrapf(err, "unable to process item %+v", item))
	}

	return true
}
