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
	"fmt"
	"os"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("OctaviaAPI controller", func() {

	var namespace string
	var octaviaAPIName types.NamespacedName

	BeforeEach(func() {
		namespace = uuid.New().String()
		th.CreateNamespace(namespace)
		DeferCleanup(th.DeleteNamespace, namespace)
		err := os.Setenv("OPERATOR_TEMPLATES", "../../templates")
		Expect(err).NotTo(HaveOccurred())

		name := fmt.Sprintf("octavia-%s", uuid.New().String())

		octaviaAPIName = CreateOctaviaAPI(namespace, name, GetDefaultOctaviaAPISpec())
		DeferCleanup(DeleteOctaviaAPI, octaviaAPIName)
	})

	When("An OctaviaAPI instance is created", func() {
		It("should have the Spec fields initialized", func() {
			OctaviaAPI := GetOctaviaAPI(octaviaAPIName)
			Expect(OctaviaAPI.Spec.DatabaseInstance).Should(Equal("test-octavia-db-instance"))
			Expect(OctaviaAPI.Spec.DatabaseUser).Should(Equal("octavia"))
			Expect(OctaviaAPI.Spec.Replicas).Should(Equal(int32(1)))
			Expect(OctaviaAPI.Spec.ServiceUser).Should(Equal("octavia"))
		})
	})

	It("should have a finalizer", func() {
		Eventually(func() []string {
			return GetOctaviaAPI(octaviaAPIName).Finalizers
		}, timeout, interval).Should(ContainElement("OctaviaAPI"))
	})

	It("should not create a config map", func() {
	})
})
