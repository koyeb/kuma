module github.com/kumahq/kuma/api

go 1.16

require (
	github.com/envoyproxy/go-control-plane v0.10.3
	github.com/envoyproxy/protoc-gen-validate v0.9.1
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/kumahq/protoc-gen-kumadoc v0.1.7
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
// When running `make generate` in this folder, one can get into errors of missing proto dependecies
// To solve the issue, uncomment the section below and run `go mod download`
//github.com/cncf/udpa latest
//github.com/envoyproxy/data-plane-api latest
//github.com/googleapis/googleapis latest
)