package play_test

import (
	"context"

	"github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("K8client CRD (Memcached)", Label(models.GINKGO_SETUP), func() {
	var (
		err       error
		namespace = "default"
		name      = "memcached-sample"
	)

	Context("using KubeConfig", func() {
		var (
			r          client.Client
			typeMeta   = metav1.TypeMeta{Kind: "Memcached", APIVersion: "cache.aman.com/v1beta1"}
			objectMeta = metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			}
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
				memcachedNew = &v1beta1.Memcached{
					TypeMeta:   typeMeta,
					ObjectMeta: objectMeta,
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
