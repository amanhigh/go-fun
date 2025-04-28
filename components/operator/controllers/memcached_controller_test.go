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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cachev1alpha1 "github.com/amanhigh/go-fun/components/operator/api/v1alpha1"
	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"github.com/amanhigh/go-fun/components/operator/common"
	"github.com/amanhigh/go-fun/models"
)

var _ = Describe("Memcached controller", Label(models.GINKGO_SETUP), func() {

	const MemcachedName = "test-memcached"

	var (
		ctx      = context.Background()
		waitTime = time.Minute
		waitStep = time.Second

		imageName    = "example.com/image:test"
		sidecarImage = common.SIDECAR_IMAGE_NAME

		namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      MemcachedName,
				Namespace: MemcachedName,
			},
		}

		memcached *cachev1beta1.Memcached
		size      = int32(1)
		port      = int32(8443)

		typeNamespaceName = types.NamespacedName{Name: MemcachedName, Namespace: MemcachedName}
		err               error
	)

	BeforeEach(func() {
		// Let's mock our custom resource at the same way that we would
		// apply on the cluster the manifest under config/samples
		memcached = &cachev1beta1.Memcached{
			ObjectMeta: metav1.ObjectMeta{
				Name:      MemcachedName,
				Namespace: namespace.Name,
			},
			Spec: cachev1beta1.MemcachedSpec{
				Size:          size,
				ContainerPort: port,
				SidecarImage:  sidecarImage,
			},
		}
	})

	Context("Memcached controller test", Ordered, func() {
		var (

			/**
				Checks if running on actual cluster as
				env test can't emulate all functions.
			**/
			isCluster = func() bool {
				_, err := os.LookupEnv("USE_EXISTING_CLUSTER")
				return err
			}
		)

		BeforeAll(func() {
			By("Creating the Namespace to perform the tests")
			Eventually(func() error {
				return k8sClient.Create(ctx, namespace)
			}, waitTime, waitStep).ShouldNot(HaveOccurred())

		})

		AfterAll(func() {
			/* Don't Delete Namespace till end due to envtest limitations. */
			// Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations. More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)
		})

		Context("Conversion", func() {
			It("should support alphav1", func() {
				memcachedAlpha1 := &cachev1alpha1.Memcached{
					ObjectMeta: memcached.ObjectMeta,
					Spec:       cachev1alpha1.MemcachedSpec{},
				}

				Expect(k8sClient.Create(ctx, memcachedAlpha1)).To(Not(HaveOccurred()))
				Expect(k8sClient.Delete(ctx, memcachedAlpha1)).To(Not(HaveOccurred()))
				// TASK: Implement Reconcile for Older Versions
			})
		})

		Context("Create Kind MemCached", func() {

			BeforeEach(func() {
				err = k8sClient.Create(ctx, memcached)
				Expect(err).To(Not(HaveOccurred()))

			})

			AfterEach(func() {
				// Clean Memcached Object if left
				_ = k8sClient.Delete(ctx, memcached)
			})

			It("should be successful", func() {
				Expect(k8sClient.Get(ctx, typeNamespaceName, memcached)).Should(Succeed())
			})

			Context("Validations", func() {
				It("should respect Max Size", func() {
					// Update Memcached CR to have size greater than Max.
					memcached.Spec.Size = int32(5)
					err = k8sClient.Update(ctx, memcached)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("spec.size"))
				})

				It("should check SidecarImage", func() {
					// Update Memcached CR to have size greater than Max.
					memcached.Spec.SidecarImage = "invalidImage"
					err = k8sClient.Update(ctx, memcached)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("supported values"))
				})

			})

			Context("Reconcile", func() {
				var (
					memcachedReconciler *MemcachedReconciler
					deployment          *appsv1.Deployment
					mgr                 ctrl.Manager
				)

				BeforeEach(func() {
					memcachedReconciler = &MemcachedReconciler{
						Client:   k8sClient,
						Scheme:   k8sClient.Scheme(),
						Recorder: record.NewFakeRecorder(10),
					}

					// Initialize helpers
					memcachedReconciler.statusHelper = NewStatusHelper(memcachedReconciler)
					memcachedReconciler.deployHelper = NewDeploymentHelper(memcachedReconciler)
					memcachedReconciler.reconcileHelper = NewReconciliationHelper(memcachedReconciler.statusHelper, memcachedReconciler.deployHelper, memcachedReconciler)

					deployment = &appsv1.Deployment{}

					mgr, err = ctrl.NewManager(cfg, ctrl.Options{
						Scheme:             k8sClient.Scheme(),
						MetricsBindAddress: "0",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not create deployment without Image ENV Variable", func() {
					By("Reconciling Creation without Image")
					_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
						NamespacedName: typeNamespaceName,
					})
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("unable to find MEMCACHED_IMAGE environment variable with the image"))

					Expect(k8sClient.Get(ctx, typeNamespaceName, deployment)).Should(HaveOccurred())

					By("Reconcile Delete for Errored Object Without Image")
					Expect(k8sClient.Delete(ctx, memcached)).To(Succeed())
					_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
						NamespacedName: typeNamespaceName,
					})
					Expect(err).ShouldNot(HaveOccurred())

					By("Reconcile In Abscense of Memcached Object shoud Do Nothing")
					_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
						NamespacedName: typeNamespaceName,
					})
					Expect(err).ShouldNot(HaveOccurred())

				})

				It("should setup", func() {
					By("Manager")
					err = memcachedReconciler.SetupWithManager(mgr)
					Expect(err).ShouldNot(HaveOccurred())

					By("Webhook")
					err = memcached.SetupWebhookWithManager(mgr)
					Expect(err).ShouldNot(HaveOccurred())
				})

				Context("Reconcile Create", func() {
					BeforeEach(func() {
						By("Setting the Image ENV VAR which stores the Operand image")
						err = os.Setenv("MEMCACHED_IMAGE", imageName)
						Expect(err).To(Not(HaveOccurred()))

						By("Reconciling Creation")
						_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
							NamespacedName: typeNamespaceName,
						})
						Expect(err).To(Not(HaveOccurred()))
					})

					AfterEach(func() {
						// Clean Memcached Object on Way Back.
						err = k8sClient.Delete(ctx, memcached)
						Expect(err).To(Not(HaveOccurred()))

						By("Reconciling Deletion")
						_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
							NamespacedName: typeNamespaceName,
						})
						Expect(err).To(Not(HaveOccurred()))

						if isCluster() {
							By("Verifying Deployment is Deleted")
							Eventually(func() error {
								return k8sClient.Get(ctx, typeNamespaceName, deployment)
							}, waitTime, waitStep).ShouldNot(Succeed())
						}

						By("Removing the Image ENV VAR which stores the Operand image")
						_ = os.Unsetenv("MEMCACHED_IMAGE")
					})

					It("should succeed for create deployment", func() {
						Eventually(func() error {
							return k8sClient.Get(ctx, typeNamespaceName, deployment)
						}, waitTime, waitStep).Should(Succeed())

						By("Verifiying Deployment Spec")
						Expect(*deployment.Spec.Replicas).To(Equal(size))
						Expect(deployment.Spec.Template.Labels).To(Equal(memcachedReconciler.deployHelper.GetLabels(memcached.Name, imageName)))
						Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(imageName))
						Expect(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsNonRoot).To(BeTrue())

						By("Verifying Sidecar")
						Expect(deployment.Spec.Template.Spec.Containers[1].Name).To(Equal("sidecar"))
						Expect(deployment.Spec.Template.Spec.Containers[1].Image).To(Equal(sidecarImage))
						Expect(deployment.Spec.Template.Spec.Containers[1].Command).To(ContainElement("sleep"))

						// Check if the Memcached object is set as the owner of the Deployment object
						Expect(deployment.ObjectMeta.OwnerReferences).To(ContainElement(metav1.OwnerReference{
							APIVersion:         "cache.aman.com/v1beta1",
							Kind:               "Memcached",
							Name:               memcached.Name,
							UID:                memcached.UID,
							Controller:         &[]bool{true}[0],
							BlockOwnerDeletion: &[]bool{true}[0],
						}))
					})

					It("should update Memcached Condition Status", func() {
						Eventually(func() error {
							if len(memcached.Status.Conditions) != 0 {
								latestStatusCondition := memcached.Status.Conditions[len(memcached.Status.Conditions)-1]
								expectedLatestStatusCondition := metav1.Condition{Type: typeAvailableMemcached,
									Status: metav1.ConditionTrue, Reason: "Reconciling",
									Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", memcached.Name, memcached.Spec.Size)}
								if latestStatusCondition != expectedLatestStatusCondition {
									return fmt.Errorf("The latest status condition added to the memcached instance is not as expected")
								}
							}
							return nil
						}, waitTime, waitStep).Should(Succeed())
					})

					Context("Scale Up", func() {
						var (
							newSize = int32(2)
						)
						BeforeEach(func() {
							// Refresh Object
							err = k8sClient.Get(ctx, typeNamespaceName, memcached)
							Expect(err).ToNot(HaveOccurred())

							// Update Memcached CR to have size of 2
							memcached.Spec.Size = newSize
							err = k8sClient.Update(ctx, memcached)
							Expect(err).To(Not(HaveOccurred()))

							By("Reconciling Scale Up")
							Eventually(func() error {
								_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
									NamespacedName: typeNamespaceName,
								})
								return err
							}, waitTime, waitStep).ShouldNot(HaveOccurred())
						})

						It("should update deployment replicas when spec size changes", func() {
							// Wait for Deployment to be updated with 2 replicas
							Eventually(func() int32 {
								err = k8sClient.Get(ctx, typeNamespaceName, deployment)
								Expect(err).ToNot(HaveOccurred())
								return *deployment.Spec.Replicas
							}, waitTime, waitStep).Should(Equal(newSize))
						})

						Context("ScaleDown", func() {
							BeforeEach(func() {
								// Reduce Cluster Size
								newSize = int32(1)

								// Refresh Object
								err = k8sClient.Get(ctx, typeNamespaceName, memcached)
								Expect(err).ToNot(HaveOccurred())

								// Update Memcached CR to have size of 1
								memcached.Spec.Size = newSize
								err = k8sClient.Update(ctx, memcached)
								Expect(err).To(Not(HaveOccurred()))

								By("Reconciling ScaleDown")
								Eventually(func() error {
									_, err = memcachedReconciler.Reconcile(ctx, reconcile.Request{
										NamespacedName: typeNamespaceName,
									})
									return err
								}, waitTime, waitStep).ShouldNot(HaveOccurred())
							})

							It("should update deployment replicas when spec size decreases", func() {
								// Wait for Deployment to be scaled down
								Eventually(func() int32 {
									err = k8sClient.Get(ctx, typeNamespaceName, deployment)
									Expect(err).ToNot(HaveOccurred())
									return *deployment.Spec.Replicas
								}, time.Minute, time.Second).Should(Equal(newSize))
							})
						})

					})

				})
			})
		})

	})

	Context("Webhook", func() {
		Context("Create Validate", func() {
			It("should succeed", func() {
				err = memcached.ValidateCreate()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should fail for wrong port", func() {
				memcached.Spec.ContainerPort = 7000
				err = memcached.ValidateCreate()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("8000"))
			})
		})

		Context("Update Validate", func() {
			var (
				oldMemcached = &cachev1alpha1.Memcached{}
			)
			It("should succeed", func() {
				err = memcached.ValidateUpdate(oldMemcached)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should fail for wrong port", func() {
				memcached.Spec.ContainerPort = 7000
				err = memcached.ValidateUpdate(oldMemcached)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("8000"))
			})
		})

		Context("Delete Validate", func() {
			It("should succeed", func() {
				err = memcached.ValidateDelete()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should fail for wrong port", func() {
				memcached.Labels = map[string]string{
					"type": "critical",
				}
				err = memcached.ValidateDelete()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("critical"))
			})
		})

		Context("Defaulting", func() {
			It("should set positive size", func() {
				memcached.Spec.Size = -1
				memcached.Default()
				Expect(memcached.Spec.Size).To(BeEquivalentTo(1))
			})
		})
	})

	It("Should fail to create without namespace", func() {
		// Attempt to create the Memcached controller
		err := k8sClient.Create(ctx, memcached)
		Expect(err).Should(HaveOccurred())
	})

})
