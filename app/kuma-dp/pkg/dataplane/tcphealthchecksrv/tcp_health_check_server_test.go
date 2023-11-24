package tcphealthchecksrv

import (
	"io"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kumahq/kuma/pkg/test"
)

var _ = Describe("TCP health check server", func() {
	c := make(chan struct{})
	srv := New("127.0.0.1", 65530, "127.0.0.1")
	go func() {
		err := srv.Start(c)
		Expect(err).NotTo(HaveOccurred())
	}()

	// Test server
	go func() {
		err := http.ListenAndServe("127.0.0.1:65531", nil)
		Expect(err).NotTo(HaveOccurred())
	}()

	// Waiting for server to boot
	time.Sleep(2 * time.Second)

	Describe("Run(..)", func() {
		It("should returns a 400 when the TCP port provided is not a number", test.Within(10*time.Second, func() {
			resp, err := http.Get("http://127.0.0.1:65530/check/helloworld")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(400))
			b, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(b)).To(Equal(`{"error":"port \"helloworld\" is invalid"}`))
		}))

		It("should returns a 400 when the TCP port provided is not a port number", test.Within(10*time.Second, func() {
			resp, err := http.Get("http://127.0.0.1:65530/check/670000")
			Expect(err).NotTo(HaveOccurred())
			b, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(400))
			Expect(string(b)).To(Equal(`{"error":"port \"670000\" is invalid"}`))
		}))

		It("should returns a 404 when the TCP port is not opened", test.Within(10*time.Second, func() {
			resp, err := http.Get("http://127.0.0.1:65530/check/1024")
			Expect(err).NotTo(HaveOccurred())
			b, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(404))
			Expect(string(b)).To(Equal(`{"error":"port \"1024\" is not listening"}`))
		}))

		It("should returns a 200 when the TCP port is not opened", test.Within(10*time.Second, func() {
			resp, err := http.Get("http://127.0.0.1:65530/check/65531")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
			b, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(b)).To(Equal(`{"message":"port \"65531\" is listening"}`))
		}))
	})
})
