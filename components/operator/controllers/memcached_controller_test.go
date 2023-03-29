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

package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cachev1alpha1 "github.com/amanhigh/go-fun/components/operator/api/v1alpha1"
	"github.com/amanhigh/go-fun/models"
)

// TODO: Include in Go Releaser
var _ = Describe("Memcached controller", Label(models.GINKGO_SETUP), func() {
	Context("Memcached controller test", Ordered, func() {

		const MemcachedName = "test-memcached"

		var (
			ctx = context.Background()

			namespace = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:      MemcachedName,
					Namespace: MemcachedName,
				},
			}
			typeNamespaceName = types.NamespacedName{Name: MemcachedName, Namespace: MemcachedName}
			err               error
		)

		BeforeAll(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))

			By("Setting the Image ENV VAR which stores the Operand image")
			err = os.Setenv("MEMCACHED_IMAGE", "example.com/image:test")
			Expect(err).To(Not(HaveOccurred()))
		})

		AfterAll(func() {
			/* Don't Delete Namespace till end due to envtest limitations. */
			// TODO(user): Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations. More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)

			By("Removing the Image ENV VAR which stores the Operand image")
			_ = os.Unsetenv("MEMCACHED_IMAGE")
		})

		Context("Create Kind MemCached", func() {
			var (
				memcached *cachev1alpha1.Memcached
				size      = int32(1)
			)

			BeforeEach(func() {
				// Let's mock our custom resource at the same way that we would
				// apply on the cluster the manifest under config/samples
				memcached = &cachev1alpha1.Memcached{
					ObjectMeta: metav1.ObjectMeta{
						Name:      MemcachedName,
						Namespace: namespace.Name,
					},
					Spec: cachev1alpha1.MemcachedSpec{
						Size: size,
					},
				}

				err = k8sClient.Create(ctx, memcached)
				Expect(err).To(Not(HaveOccurred()))

			})

			AfterEach(func() {
				//Clean Memcached Object if left
				_ = k8sClient.Delete(ctx, memcached)
			})

			It("should be successful", func() {
				Eventually(func() error {
					found := &cachev1alpha1.Memcached{}
					return k8sClient.Get(ctx, typeNamespaceName, found)
				}, time.Minute, time.Second).Should(Succeed())
			})

			Context("Reconcile", func() {
				var (
					memcachedReconciler *MemcachedReconciler
				)

				BeforeEach(func() {
					memcachedReconciler = &MemcachedReconciler{
						Client:   k8sClient,
						Scheme:   k8sClient.Scheme(),
						Recorder: record.NewFakeRecorder(10),
					}

					By("Reconciling Creation")
					_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
						NamespacedName: typeNamespaceName,
					})
					Expect(err).To(Not(HaveOccurred()))
				})

				AfterEach(func() {
					//Clean Memcached Object on Way Back.
					err = k8sClient.Delete(ctx, memcached)
					Expect(err).To(Not(HaveOccurred()))

					By("Reconciling Deletion")
					_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
						NamespacedName: typeNamespaceName,
					})
					Expect(err).To(Not(HaveOccurred()))
				})

				It("should succeed for create deployment", func() {
					By("Checking if Deployment was successfully created in the reconciliation")
					Eventually(func() error {
						found := &appsv1.Deployment{}
						return k8sClient.Get(ctx, typeNamespaceName, found)
					}, time.Minute, time.Second).Should(Succeed())
				})

				It("should update Memcached Condition Status", func() {
					Eventually(func() error {
						if memcached.Status.Conditions != nil && len(memcached.Status.Conditions) != 0 {
							latestStatusCondition := memcached.Status.Conditions[len(memcached.Status.Conditions)-1]
							expectedLatestStatusCondition := metav1.Condition{Type: typeAvailableMemcached,
								Status: metav1.ConditionTrue, Reason: "Reconciling",
								Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", memcached.Name, memcached.Spec.Size)}
							if latestStatusCondition != expectedLatestStatusCondition {
								return fmt.Errorf("The latest status condition added to the memcached instance is not as expected")
							}
						}
						return nil
					}, time.Minute, time.Second).Should(Succeed())
				})
			})
		})
	})
})
