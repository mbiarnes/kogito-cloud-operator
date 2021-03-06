// Copyright 2019 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	kafkabetav1 "github.com/kiegroup/kogito-cloud-operator/pkg/apis/kafka/v1beta1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure"
	"github.com/kiegroup/kogito-cloud-operator/pkg/logger"
	v1 "github.com/openshift/api/image/v1"

	routev1 "github.com/openshift/api/route/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var log = logger.GetLogger("resources_kogitodataindex")

// KogitoDataIndexResources is the data structure for all resources managed by KogitoDataIndex CR
type KogitoDataIndexResources struct {
	// Deployment is the resource responsible for the Data Index Service image deployment in the cluster
	Deployment *appsv1.Deployment
	// Service to expose the data index service internally
	Service *corev1.Service
	// Route to expose the service in the Ingress. Supported only on OpenShift for now
	Route *routev1.Route
	// KafkaTopics are the Kafka Topics required by the Data Index Service
	KafkaTopics []*kafkabetav1.KafkaTopic
	// ImageStream is the image stream OpenShift resource to manage images (not available on k8s)
	ImageStream *v1.ImageStream
}

// KogitoDataIndexResourcesStatus identifies the status of the resource
type KogitoDataIndexResourcesStatus struct {
	New bool
}

type kogitoDataIndexResourcesFactory struct {
	Client          *client.Client
	Resources       *KogitoDataIndexResources
	KogitoDataIndex *v1alpha1.KogitoDataIndex
	Error           error
}

// Build will call a builder function if no errors were found
func (f *kogitoDataIndexResourcesFactory) build(fn func(*kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	if f.Error == nil {
		return fn(f)
	}
	// break the chain
	return f
}

func (f *kogitoDataIndexResourcesFactory) buildOnOpenshift(fn func(*kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	if f.Error == nil && f.Client.IsOpenshift() {
		return fn(f)
	}
	// break the chain
	return f
}

// GetRequestedResources will create the data of the needed resources objects for Kogito Data Index service
func GetRequestedResources(instance *v1alpha1.KogitoDataIndex, client *client.Client) (*KogitoDataIndexResources, error) {
	factory := kogitoDataIndexResourcesFactory{
		Client:          client,
		Resources:       &KogitoDataIndexResources{},
		KogitoDataIndex: instance,
	}

	factory.
		buildOnOpenshift(createImageStream).
		build(createDeployment).
		build(createService).
		buildOnOpenshift(createRoute).
		build(createKafkaTopic)

	return factory.Resources, factory.Error
}

func createDeployment(f *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	infinispanCredentials, err := infrastructure.FetchInfinispanCredentials(&f.KogitoDataIndex.Spec, f.KogitoDataIndex.Namespace, f.Client)
	if err != nil {
		f.Error = err
		return f
	}
	kafkaServerURI, err := getKafkaServerURI(f.KogitoDataIndex.Spec.KafkaProperties, f.KogitoDataIndex.Namespace, f.Client)
	if err != nil {
		f.Error = err
		return f
	}
	deployment, err := newDeployment(f.KogitoDataIndex, infinispanCredentials, kafkaServerURI, f.Client, f.Resources.ImageStream)
	if err != nil {
		f.Error = err
		return f
	}
	f.Resources.Deployment = deployment
	return f
}

func createService(f *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	svc := newService(f.KogitoDataIndex, f.Resources.Deployment)
	f.Resources.Service = svc
	return f
}

func createRoute(f *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	route, err := newRoute(f.KogitoDataIndex, f.Resources.Service)
	if err != nil {
		f.Error = err
		return f
	}
	f.Resources.Route = route
	return f
}

func createKafkaTopic(f *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	kafkaName, kafkaReplicas, err := getKafkaServerReplicas(f.KogitoDataIndex.Spec.KafkaProperties, f.KogitoDataIndex.Namespace, f.Client)
	if err != nil {
		f.Error = err
		return f
	} else if len(kafkaName) <= 0 || kafkaReplicas <= 0 {
		return f
	}

	for _, kafkaTopicName := range kafkaTopicNames {
		kafkaTopic := newKafkaTopic(kafkaTopicName, kafkaName, kafkaReplicas, f.KogitoDataIndex.Namespace)
		f.Resources.KafkaTopics = append(f.Resources.KafkaTopics, kafkaTopic)
	}

	return f
}

func createImageStream(f *kogitoDataIndexResourcesFactory) *kogitoDataIndexResourcesFactory {
	f.Resources.ImageStream = newImage(f.KogitoDataIndex)
	return f
}
