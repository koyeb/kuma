package tcphealthchecksrv

import (
	"testing"

	"github.com/kumahq/kuma/pkg/test"
)

func TestTCPHealthCheckServer(t *testing.T) {
	test.RunSpecs(t, "TCP health check server")
}
