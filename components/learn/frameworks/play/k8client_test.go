package play_test

import (
	"context"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var _ = Describe("K8client", Label(models.GINKGO_SETUP), func() {
	var (
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		err        error
		clientset  *kubernetes.Clientset
		waitTime   = time.Second * 30
		namespace  = "default"
	)

	Context("using Config", func() {
		BeforeEach(func() {
			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			Expect(err).ShouldNot(HaveOccurred())

			clientset, err = kubernetes.NewForConfig(config)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should build", func() {
			Expect(clientset).ShouldNot(BeNil())
		})

		Context("Deployment Create", func() {
			var (
				size              = int32(1)
				deploymentName    = "mysql-deployment"
				deploymentsClient v1.DeploymentInterface
			)

			//Spec Vars
			var (
				// Define Selector
				selector = &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "mysql",
					},
				}

				//Object Meta
				objectMeta = metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mysql",
					},
				}

				//Pod Spec
				podspec = corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql",
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "root",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
						},
					},
				}

				// Define Template
				template = corev1.PodTemplateSpec{
					ObjectMeta: objectMeta,
					Spec:       podspec,
				}

				// Define Deployment
				deployment = &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: deploymentName,
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &size,
						Selector: selector,
						Template: template,
					},
				}
			)

			BeforeEach(func() {
				deploymentsClient = clientset.AppsV1().Deployments(namespace)
				By("Creating MySQL deployment...")
				_, err = deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

			})

			AfterEach(func() {
				By("Deleting MySQL deployment...")
				err = deploymentsClient.Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have correct status", func() {
				Eventually(func() (size int32) {
					deployment, err := deploymentsClient.Get(context.Background(), deploymentName, metav1.GetOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					return *deployment.Spec.Replicas
				}, waitTime).Should(Equal(size))
			})

			Context("Resize Deployment", func() {
				var (
					newSize = size + 1
				)
				BeforeEach(func() {
					By("Resizing MySQL deployment...")
					// Modify the deployment to set the replicas to 2
					deployment.Spec.Replicas = &newSize

					// Create the deployment with 2 replicas
					_, err = deploymentsClient.Update(context.Background(), deployment, metav1.UpdateOptions{})

					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should work", func() {
					// Wait for the deployment to reach the expected size
					Eventually(func() int32 {
						deployment, err := deploymentsClient.Get(context.Background(), deploymentName, metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())
						return deployment.Status.Replicas
					}, waitTime).Should(Equal(newSize))

				})

			})
		})
	})

	Context("using KubeConfig", func() {
		var (
			r client.Client
		)
		BeforeEach(func() {
			config := config.GetConfigOrDie()
			v1beta1.AddToScheme(scheme.Scheme)
			r, err = client.New(config, client.Options{Scheme: scheme.Scheme})
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should build", func() {
			Expect(r).ShouldNot(BeNil())
		})
		Context("Create", func() {
			var (
				name = "memcached-sample"

				memcachedNew = &v1beta1.Memcached{
					TypeMeta: metav1.TypeMeta{Kind: "Memcached", APIVersion: "cache.aman.com/v1beta1"},
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: v1beta1.MemcachedSpec{
						Size:          2,
						ContainerPort: 8443,
						SidecarImage:  "busybox",
					},
				}
			)
			BeforeEach(func() {
				err = r.Create(context.Background(), memcachedNew)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err = r.Delete(context.Background(), memcachedNew)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should get Memcached", func() {
				memcached := &v1beta1.Memcached{}
				err = r.Get(context.Background(), client.ObjectKey{
					Namespace: namespace,
					Name:      name,
				}, memcached)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
