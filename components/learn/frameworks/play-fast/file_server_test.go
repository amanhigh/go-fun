package play_fast_test

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Server", func() {
	var (
		// create file server handler
		dir  = "/tmp"
		fs   = http.FileServer(http.Dir(dir))
		port = 8092
		url  = "http://localhost:" + strconv.Itoa(port)

		response *resty.Response
		err      error
	)
	BeforeEach(func() {
		srv := util.NewTestServer(fmt.Sprintf(":%v", port))
		srv.Handler = fs
		go srv.ListenAndServe() //nolint:errcheck
	})

	It("should run", func() {
		response, err = resty.New().R().Get(url)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode()).To(Equal(http.StatusOK))
		Expect(response.String()).To(Not(BeNil()))
	})

})
