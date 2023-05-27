package play_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var _ = FDescribe("K8client", func() {
	var (
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		// err        error
		clientset *kubernetes.Clientset
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
	})
})
