/*
Copyright 2023.

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

package functional_test

import (
	"time"

	. "github.com/onsi/gomega"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	octaviav1 "github.com/openstack-k8s-operators/octavia-operator/api/v1beta1"
)

const (
	SecretName = "octavia-test-secret"
	timeout    = time.Second * 10
	interval   = timeout / 100
)

func GetDefaultOctaviaAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance": "test-octavia-db-instance",
		"containerImage":   "test-octavia-api-image",
		"secret":           SecretName,
	}
}

func CreateOctaviaAPI(namespace string, OctaviaAPIName string, spec map[string]interface{}) types.NamespacedName {
	raw := map[string]interface{}{
		"apiVersion": "octavia.openstack.org/v1beta1",
		"kind":       "OctaviaAPI",
		"metadata": map[string]interface{}{
			"name":      OctaviaAPIName,
			"namespace": namespace,
		},
		"spec": spec,
	}
	th.CreateUnstructured(raw)

	return types.NamespacedName{Name: OctaviaAPIName, Namespace: namespace}
}

func DeleteOctaviaAPI(name types.NamespacedName) {
	Eventually(func(g Gomega) {
		OctaviaAPI := &octaviav1.OctaviaAPI{}
		err := k8sClient.Get(ctx, name, OctaviaAPI)
		if k8s_errors.IsNotFound(err) {
			// No need to do more!
			return
		}
		g.Expect(err).Should(BeNil())
		g.Expect(k8sClient.Delete(ctx, OctaviaAPI)).Should(Succeed())

		err = k8sClient.Get(ctx, name, OctaviaAPI)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
}

func GetOctaviaAPI(name types.NamespacedName) *octaviav1.OctaviaAPI {
	instance := &octaviav1.OctaviaAPI{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func OctaviaAPIConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetOctaviaAPI(name)
	return instance.Status.Conditions
}
