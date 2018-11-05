/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"time"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"

	"github.com/golang/glog"
	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	//appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	samplev1alpha1 "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/apis/cryptopussecretcontroller/v1alpha1"
	clientset "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/client/clientset/versioned"
	samplescheme "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/client/clientset/versioned/scheme"
	informers "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/client/informers/externalversions/cryptopussecretcontroller/v1alpha1"
	listers "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/client/listers/cryptopussecretcontroller/v1alpha1"
)

const controllerAgentName = "cryptopus-k8s-secretcontroller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a SecretClaim is synced
	SuccessSynced = "Synced"
	FailedSynced = "Failed"

	// ErrResourceExists is used as part of the Event 'reason' when a SecretClaim fails
	// to sync due to a Secret of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Secret already existing
	MessageResourceExists = "Resource %q already exists and is not managed by SecretClaim"


	// MessageResourceSynced is the message used for an Event fired when a SecretClaim
	// is synced successfully
	MessageResourceSynced = "SecretClaim synced successfully"

	MessageCryptopusAPIConfigNotFound = "Secret '%s' with Cryptopus API Details not found"
	MessageCryptopusAPIInvalid = "Cryptopus API Config Invalid. Make sure secret '%s' contains CRYPTOPUS_API, CRYPTOPUS_API_USER, CRYPTOPUS_API_TOKEN"
	MessageCryptopusAPIFailed = "Unable to make request to Cryptopus API: '%s': %s"
	MessageCryptopusAPINewRequestFailed = "Unable to create HTTP Request Client for Cryptopus API: '%s': %s"

	MessageFailed = "Unable to get Secret for Account %s sfrom Cryptopus: %s"

	MessageSecretClaimValueRequired = "The value '%s' is required in Secretclaim '%s'"


)


type CryptopusAccountDetail struct{
    Id      int32 `json:"id"`
    AccountName string `json:"accountname"`
		GroupId int32 `json:"group_id"`
		Group string `json:"group"`
		Team string `json:"team"`
		TeamId int32 `json:"team_id"`
		CleartextPassword string `json:"cleartext_password"`
		CleartextUsername string `json:"cleartext_username"`
}

type CryptopusAccountList struct{
    Account CryptopusAccountDetail `json:"account"`
}

type CryptopusMessage struct{
    errors      []string `json:"errors"`
    info 				[]string `json:"info"`
}

type CryptopusResponse struct{
    Data      CryptopusAccountList `json:"data"`
    Message CryptopusMessage `json:"messages"`
}



// Controller is the controller implementation for SecretClaim resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// vaultcontrollerclientset is a clientset for our own API group
	cryptopuscontrollerclientset clientset.Interface
	secretsLister 					corelisters.SecretLister
	secretsSynced 					cache.InformerSynced
	secretClaimsLister   	listers.SecretClaimLister
	secretClaimsSynced      cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
	defaultRefreshTime int32
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	cryptopuscontrollerclientset clientset.Interface,
	secretInformer coreinformers.SecretInformer,
	secretClaimInformer informers.SecretClaimInformer,
	defaultRefreshTime int32) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:     kubeclientset,
		cryptopuscontrollerclientset:   cryptopuscontrollerclientset,

		secretsLister: secretInformer.Lister(),
		secretsSynced: secretInformer.Informer().HasSynced,

		secretClaimsLister:        secretClaimInformer.Lister(),
		secretClaimsSynced:        secretClaimInformer.Informer().HasSynced,

		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "SecretClaims"),
		recorder:          recorder,
		defaultRefreshTime: defaultRefreshTime,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	secretClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueSecretClaim,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueSecretClaim(new)
		},
		DeleteFunc: controller.enqueueSecretClaim,
	})
	// Set up an event handler for when Secret resources change. This
	// handler will lookup the owner of the given Secret, and if it is
	// owned by a SecretClaim resource will enqueue that SecretClaim resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Secret resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	secretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newSecret := new.(*corev1.Secret)
			oldSecret := old.(*corev1.Secret)
			if newSecret.ResourceVersion == oldSecret.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Secret will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
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
	glog.Info("Starting SecretClaim controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.secretsSynced, c.secretClaimsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process SecretClaim resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
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
		// SecretClaim resource to be synced.
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
// converge the two. It then updates the Status block of the SecretClaim resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the SecretClaim resource with this namespace/name
	secretClaim, err := c.secretClaimsLister.SecretClaims(namespace).Get(name)
	if err != nil {
		// The SecretClaim resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("SecretClaim '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	err = valueRequired(c, secretClaim, "secretName", secretClaim.Spec.SecretName)
	if err != nil {
		return err
	}


	// Get the Secret with the secretName specified in SecretClaim.spec
	secret, err := c.secretsLister.Secrets(secretClaim.Namespace).Get(secretClaim.Spec.SecretName)


	// If the Secret is not controlled by this SecretClaim resource, we should log
	// a warning to the event recorder and ret
	if secret != nil && !metav1.IsControlledBy(secret, secretClaim) {
		msg := fmt.Sprintf(MessageResourceExists, secret.Name)
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// Finally, we update the status block of the SecretClaim resource to reflect the
	// current state of the world
	err = c.updateSecretClaimStatus(secretClaim, secret)
	if err != nil {
		return err
	}

	return nil
}


// Make sure value is not ""
func valueRequired(c *Controller, secretClaim *samplev1alpha1.SecretClaim, key string, value string) error {
	res := value
	if res == "" {

		msg := fmt.Sprintf(MessageSecretClaimValueRequired, key, secretClaim.Name)
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
		return fmt.Errorf(msg)

	}
	return nil
}

func (c *Controller) updateSecretClaimStatus(secretClaim *samplev1alpha1.SecretClaim, secret *corev1.Secret) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	secretClaimCopy := secretClaim.DeepCopy()
	secretClaimCopy.Status.Phase = "Bound"

	// Calculate diff in seconds since last Update
	now := time.Now()
	diff := int32(now.Sub(time.Unix(int64(secretClaimCopy.Status.LastUpdate), 0)).Seconds())

	// Use Default Update time (ENV CONTROLLER_DEFAULT_REFRESH_TIME or 300 s) if not specified in SecretClaim
	updateTime := c.defaultRefreshTime
	if secretClaim.Spec.RefreshTime > 0 {
		updateTime = secretClaim.Spec.RefreshTime
	}

	if diff > updateTime || secret == nil {
		glog.V(4).Infof("Updating Secret for SecretClaim: '%s'. Last Update: '%d'",secretClaimCopy.Name, secretClaimCopy.Status.LastUpdate)

		// Update the Secret
		if updateSecret(c, secretClaimCopy, secret) {
			secretClaimCopy.Status.LastUpdate = int32(now.Unix())
			_, err := c.cryptopuscontrollerclientset.CryptopussecretcontrollerV1alpha1().SecretClaims(secretClaim.Namespace).Update(secretClaimCopy)
			if err == nil {
				c.recorder.Event(secretClaim, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
			}
		}
	}

	return nil
}

// enqueueSecretClaim takes a SecretClaim resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than SecretClaim.
func (c *Controller) enqueueSecretClaim(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the SecretClaim resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that SecretClaim resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.V(4).Infof("Processing object: %s", object.GetName())

	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a SecretClaim, we should not do anything more
		// with it.
		if ownerRef.Kind != "SecretClaim" {
			return
		}

		// Get the Owner (= SecretClaim)
		secretClaim, err := c.secretClaimsLister.SecretClaims(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.V(4).Infof("SecretClaim with Name '%s' does not exists (anymore)", ownerRef.Name)
			return
		}

		// Enque for Update Handling
		c.enqueueSecretClaim(secretClaim)
		return
	}
}

// Update a Secret with Values from Vault
func updateSecret(c *Controller, secretClaim *samplev1alpha1.SecretClaim, secret *corev1.Secret) bool {

	secretNamespace, secretName, err := cache.SplitMetaNamespaceKey(secretClaim.Spec.CryptopusSecret)

	cryptopusSecret, err := c.secretsLister.Secrets(secretNamespace).Get(secretName)
	if errors.IsNotFound(err) {
		msg := fmt.Sprintf(MessageCryptopusAPIConfigNotFound, secretName)
		glog.V(4).Infof(msg)
		secretClaim.Status.Phase = "Error"
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
		return false
	}

	cryptopus_api := string(cryptopusSecret.Data["CRYPTOPUS_API"])
	cryptopus_api_user := string(cryptopusSecret.Data["CRYPTOPUS_API_USER"])
	cryptopus_api_token := base64.StdEncoding.EncodeToString(cryptopusSecret.Data["CRYPTOPUS_API_TOKEN"])

	if len(cryptopus_api) == 0 || len(cryptopus_api_user) == 0 || len(cryptopus_api_token) == 0 {
		msg := fmt.Sprintf(MessageCryptopusAPIInvalid, secretName)
		glog.V(4).Infof(msg)
		secretClaim.Status.Phase = "Error"
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
		return false
	}

	var cryptopus_responses []CryptopusResponse

	for i := range secretClaim.Spec.Id {
		cryptopus_response := getSecret(c, secretClaim,  secretClaim.Spec.Id[i], cryptopus_api, cryptopus_api_user, cryptopus_api_token )
		if cryptopus_response == nil {
			continue
		}
		cryptopus_responses = append(cryptopus_responses, *cryptopus_response )
	}


	// Get the Secret with the secretName specified in SecretClaim.spec
	// If the resource doesn't exist, we'll create it
	secret, err = c.secretsLister.Secrets(secretClaim.Namespace).Get(secretClaim.Spec.SecretName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
	  secret, err = c.kubeclientset.CoreV1().Secrets(secretClaim.Namespace).Create(newSecret(c, secretClaim))

		if err != nil {
			msg := fmt.Sprintf("Failed to create secret: %s", err)
			glog.V(4).Infof(msg)
		}
	}
	// copy to work with
	secretcopy := secret.DeepCopy()


	// Update Secret
	annotations := map[string]string{}
	secretcopy.ObjectMeta.Annotations = annotations
	data := make(map[string][]byte, len(cryptopus_responses))
	for i := range cryptopus_responses {
		data[fmt.Sprintf("username_%d", cryptopus_responses[i].Data.Account.Id)] = []byte(cryptopus_responses[i].Data.Account.CleartextUsername)
		data[fmt.Sprintf("password_%d", cryptopus_responses[i].Data.Account.Id)] = []byte(cryptopus_responses[i].Data.Account.CleartextPassword)
	}

	secretcopy.Data = data


	_, err = c.kubeclientset.CoreV1().Secrets(secretClaim.Namespace).Update(secretcopy)
	if err != nil {
		glog.V(4).Infof("Failed to update secret '%s'", secret)
		return false
	}

	return true
}

func getSecret(c *Controller, secretClaim *samplev1alpha1.SecretClaim, cryptopus_account_id int32, cryptopus_api string, cryptopus_api_user string, cryptopus_api_token string ) *CryptopusResponse {


	tr := &http.Transport{
		// TODO: fore some Reason, we cannot verify the Certificate in our OpenShift Cluster, when running the Controller in Cluster.. I have to investigate this.
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s/api/accounts/%d}", cryptopus_api, cryptopus_account_id)
	msg := fmt.Sprintf("API: %s, URL: %s, User: %s", cryptopus_api, url, cryptopus_api_user)
	glog.V(4).Infof(msg)

	req, err := http.NewRequest("GET", url , nil)

	if err != nil {
		msg := fmt.Sprintf(MessageCryptopusAPINewRequestFailed, cryptopus_api, err)
		glog.V(4).Infof(msg)
		secretClaim.Status.Phase = "Error"
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
		return nil
	}

	// Authenticatoin Header
	req.Header.Add("Authorization-User", cryptopus_api_user)
	req.Header.Add("Authorization-Password", cryptopus_api_token)
	resp, err := client.Do(req)


	if err != nil {
			msg := fmt.Sprintf(MessageCryptopusAPIFailed, cryptopus_api, err)
			glog.V(4).Infof(msg)
			secretClaim.Status.Phase = "Error"
			c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
			return nil
		}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf(MessageFailed, cryptopus_account_id, body)
		glog.V(4).Infof(msg)
		secretClaim.Status.Phase = "Error"
		c.recorder.Event(secretClaim, corev1.EventTypeWarning, FailedSynced, msg)
		return nil
	}


	var cryptopus_response CryptopusResponse
	err = json.Unmarshal(body, &cryptopus_response)

	return &cryptopus_response
}


// newSecret creates a new Secret for a SecretClaim resource. It also sets
// the appropriate lable on the resource so handleObject can discover
// the SecretClaim resource that 'owns' it.
func newSecret(c *Controller, secretClaim *samplev1alpha1.SecretClaim) *corev1.Secret {

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretClaim.Spec.SecretName,
			Namespace: secretClaim.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(secretClaim, schema.GroupVersionKind{
					Group:   samplev1alpha1.SchemeGroupVersion.Group,
					Version: samplev1alpha1.SchemeGroupVersion.Version,
					Kind:    "SecretClaim",
				}),
			},
		},
	}
}
