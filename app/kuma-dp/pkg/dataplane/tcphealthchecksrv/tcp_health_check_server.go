package tcphealthchecksrv

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kumahq/kuma/pkg/core"
)

var (
	runLog = core.Log.WithName("tcp-health-check-server")
)

// We use the maximal TCP timeout allowed and add a 5s jitter to avoid any race condition
const defaultTcpConnectionTimeout = time.Minute*10 + time.Second*5

// TCPHealthCheckServer provides an HTTP server that will
// perform a TCP connection from a port provided as parameter
type TCPHealthCheckServer struct {
	serverIP   string
	serverPort uint32
	tcpProbeIP string
}

func New(serverIP string, serverPort uint32, tcpProbeIP string) TCPHealthCheckServer {
	return TCPHealthCheckServer{
		serverIP:   serverIP,
		serverPort: serverPort,
		tcpProbeIP: tcpProbeIP,
	}
}

func (t TCPHealthCheckServer) NeedLeaderElection() bool {
	return false
}

func (t TCPHealthCheckServer) registerHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/check/", func(w http.ResponseWriter, r *http.Request) {
		portStr := strings.TrimPrefix(r.URL.Path, "/check/")
		w.Header().Set("Content-Type", "application/json")

		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			runLog.Error(err, "port is invalid", "port", portStr)
			b, err := json.Marshal(map[string]string{"error": fmt.Sprintf(`port "%s" is invalid`, portStr)})
			if err != nil {
				runLog.Error(err, "failure when serializing JSON response")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(b)
			if err != nil {
				runLog.Error(err, "failure when writing HTTP response body")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		d := net.Dialer{Timeout: defaultTcpConnectionTimeout}
		conn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", t.tcpProbeIP, port))
		if err != nil {
			b, err := json.Marshal(map[string]string{"error": fmt.Sprintf(`port "%d" is not listening`, port)})
			if err != nil {
				runLog.Error(err, "failure when serializing JSON response")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write(b)
			if err != nil {
				runLog.Error(err, "failure when writing HTTP response body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			b, err := json.Marshal(map[string]string{"message": fmt.Sprintf(`port "%d" is listening`, port)})
			if err != nil {
				runLog.Error(err, "failure when serializing JSON response")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(b)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				runLog.Error(err, "failure when writing HTTP response body")
				return
			}
			err = conn.Close()
			if err != nil {
				runLog.Error(err, "failure when closing TCP connection")
			}
		}
	})
	return mux
}

func (t TCPHealthCheckServer) Start(stop <-chan struct{}) error {
	errCh := make(chan error, 1)
	go (func() {
		runLog.Info("starting the TCP health check server", "addr", t.serverIP, "port", t.serverPort)
		mux := t.registerHandler()
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", t.serverIP, t.serverPort), mux)
		if err != nil {
			errCh <- err
		}
	})()
	select {
	case err := <-errCh:
		return err
	case <-stop:
		runLog.Info("stopping the TCP health check server")
		return nil
	}
}
