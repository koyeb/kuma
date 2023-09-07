
set dotenv-load

use_debugger := if env_var_or_default("DEBUGGER", "") != "" { "yes" } else { "no" }
artifacts := "./build/artifacts-darwin-amd64"
kumactl := artifacts / "kumactl/kumactl"
kuma-cp := artifacts / "kuma-cp/kuma-cp"
kuma-dp := artifacts / "kuma-dp/kuma-dp"

dev-kumactl := if use_debugger == "yes" { "dlv debug app/kumactl/main.go --" } else { kumactl }
dev-kuma-cp := if use_debugger == "yes" { "dlv debug app/kuma-cp/main.go --" } else { kuma-cp }
dev-kuma-dp := if use_debugger == "yes" { "dlv debug app/kuma-dp/main.go --" } else { kuma-dp }

certs-path := "./build/koyeb/certs"

default:
  @just --list

run-db:
  -docker run -p 5532:5432 --rm --name kuma-db -e POSTGRES_USER=kuma -e POSTGRES_PASSWORD=kuma -e POSTGRES_DB=kuma -d postgres:14
  sleep 5
  docker run --rm --link kuma-db:kuma-db -e PGHOST=kuma-db -e PGUSER=kuma -e PGPASSWORD=kuma postgres:14 /bin/bash -c "echo 'CREATE DATABASE \"cp-global\" ENCODING UTF8;' | psql"
  docker run --rm --link kuma-db:kuma-db -e PGHOST=kuma-db -e PGUSER=kuma -e PGPASSWORD=kuma postgres:14 /bin/bash -c "echo 'CREATE DATABASE \"cp-par1\" ENCODING UTF8;' | psql"

stop-db:
  docker stop kuma-db

# This generates a CA cert and its key. Those are expected to be loaded into each mesh
generate-ca-cert:
  mkdir -p {{certs-path}}
  echo "\n[req]\ndistinguished_name=dn\n[ dn ]\n[ ext ]\nbasicConstraints=CA:TRUE,pathlen:0\nkeyUsage=keyCertSign\n" > /tmp/ca_config
  openssl req -config /tmp/ca_config -new -newkey rsa:2048 -nodes -subj "/CN=Hello" -x509 -extensions ext -keyout {{certs-path}}/key.pem -out {{certs-path}}/crt.pem

_build-cp:
  make build/kuma-cp

_build-dp:
  make build/kuma-dp
  make build/kumactl

cp-global: _build-cp
  {{kuma-cp}} -c koyeb/samples/cp-global.yaml migrate up
  {{dev-kuma-cp}} -c koyeb/samples/cp-global.yaml run

cp-par1: _build-cp
  {{kuma-cp}} -c koyeb/samples/cp-par1.yaml migrate up
  {{dev-kuma-cp}} -c koyeb/samples/cp-par1.yaml run

_create_secret name mesh path:
  export SECRET_NAME={{name}} ; \
  export SECRET_MESH={{mesh}} ; \
  export SECRET_BASE64_ENCODED=$(cat {{path}} | base64) ; \
  cat koyeb/samples/secret-template.yaml | envsubst | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -

_inject_ca mesh:
  @just _create_secret manually-generated-ca-cert {{mesh}} ./build/koyeb/certs/crt.pem
  @just _create_secret manually-generated-ca-key {{mesh}} ./build/koyeb/certs/key.pem

_init-default: (_inject_ca "default")
  # Upsert default mesh
  cat koyeb/samples/default-mesh.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  cat koyeb/samples/default-meshtrace.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -

ingress: _build-dp _init-default
  {{kumactl}} generate zone-token --zone=par1 --valid-for 720h --scope ingress > /tmp/dp-token-ingress
  {{dev-kuma-dp}} run \
    --proxy-type=ingress \
    --cp-address=https://127.0.0.1:5678 \
    --dataplane-token-file=/tmp/dp-token-ingress \
    --dataplane-file=koyeb/samples/ingress-par1.yaml


#glb: gen-certs _build-dp _init-default
#  {{artifacts}}/kumactl/kumactl generate zone-ingress-token --zone par1  > /tmp/dp-token-glb
#  {{artifacts}}/kuma-dp/kuma-dp run --dataplane-token-file /tmp/dp-token-glb --dns-coredns-config-template-path ./koyeb/samples/Corefile --dns-coredns-port 10053 --dns-envoy-port 10050 --log-level info --cp-address https://localhost:5678 --ca-cert-file ./build/koyeb/tls-cert/ca.pem --admin-port 4243 -d ./koyeb/samples/ingress-glb.yaml  --proxy-type ingress

igw: _build-dp _init-default
  cat koyeb/samples/default-meshgateway.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -

  {{kumactl}} generate dataplane-token -m default --valid-for 720h --config-file koyeb/samples/kumactl-configs/par1-cp.yaml > /tmp/igw-par1-token
  {{dev-kuma-dp}} run --dataplane-token-file /tmp/igw-par1-token --log-level info --cp-address https://localhost:5678 -d ./koyeb/samples/ingress-gateway-par1.yaml

#test-grpc:
#  grpcurl --plaintext localhost:8004 main.HelloWorld/Greeting
#  # evans required for advance tls config for igw
#  #echo '{}' | evans --verbose --host localhost --port 5601 --tls --certkey ./build/koyeb/gateway-client/cert-key.pem --cert ./build/koyeb/gateway-client/cert.pem --cacert ./build/koyeb/gateway-client/ca.pem -r cli --header "x-koyeb-route=dp-8004_prod" call main.HelloWorld.Greeting
#  grpcurl -authority grpc.local.koyeb.app -plaintext -servername grpc.local.koyeb.app localhost:5600 list
#  grpcurl -authority grpc.local.koyeb.app -plaintext -servername grpc.local.koyeb.app localhost:5600 main.HelloWorld/Greeting
#
#test-ws-websocat:
#  echo "hello" | websocat -t - ws-ll-c:http-request:tcp:127.0.0.1:5600  \
#    --request-header 'Host: http.local.koyeb.app' \
#    --request-header 'Upgrade: websocket' \
#    --request-header 'Sec-WebSocket-Key: mYUkMl6bemnLatx/g7ySfw==' \
#    --request-header 'Sec-WebSocket-Version: 13' \
#    --request-header 'Connection: Upgrade' \
#    --request-uri=/ -1
#
#test-ws-python:
#  python koyeb/samples/websocket-client.py
#
#test-http2:
#  curl http://localhost:8002/health -s --fail --output /dev/null
#  curl http://localhost:8002/health --http2 -s --fail  --output /dev/null
#  curl http://localhost:8002/health --http2-prior-knowledge -s --fail --output /dev/null
#  curl http://localhost:5600/health -H "host: http2.local.koyeb.app" -s --fail --output /dev/null
#  curl http://localhost:5600/health -H "host: http2.local.koyeb.app" --http2 -s --fail --output /dev/null
#  curl http://localhost:5600/health -H "host: http2.local.koyeb.app" --http2-prior-knowledge -s --fail --output /dev/null
#
#test-http:
#  curl http://localhost:8001/health -s --fail --output /dev/null
#  curl http://localhost:5600/health
#  @echo
#  curl http://127.0.0.1:5600/health
#  @echo
#  curl https://localhost:5601/health --key ./build/koyeb/gateway-cert/gateway-key.pem --cert ./build/koyeb/gateway-cert/gateway.pem --cacert ./build/koyeb/gateway-cert/ca.pem
#  @echo
#  curl https://127.0.0.1:5601/health --key ./build/koyeb/gateway-cert/gateway-key.pem --cert ./build/koyeb/gateway-cert/gateway.pem --cacert ./build/koyeb/gateway-cert/ca.pem
#  @echo
#  curl https://localhost:5601/health --key ./build/koyeb/gateway-client/cert-key.pem --cert ./build/koyeb/gateway-client/cert.pem --cacert ./build/koyeb/gateway-client/ca.pem
#  @echo
#  curl https://127.0.0.1:5601/health --key ./build/koyeb/gateway-client/cert-key.pem --cert ./build/koyeb/gateway-client/cert.pem --cacert ./build/koyeb/gateway-client/ca.pem
#  @echo
#  curl https://localhost:5601 --key ./build/koyeb/gateway-cert/tls-key.pem --cert ./build/koyeb/gateway-cert/tls.pem --cacert ./build/koyeb/gateway-cert/ca.pem -H "x-koyeb-route: dp-8001_prod" -s --fail --output /dev/null
#  curl http://localhost:5600 -H "host: http.local.koyeb.app" -s --fail --output /dev/null
#  curl http://localhost:5600 -H "host: http.local.koyeb.app" --http2-prior-knowledge -s --fail --output /dev/null
#
#test: test-http test-http2 test-grpc
#
#dp-container-ws:
#  docker run -ti -p 8001:8080 jmalloc/echo-server

dp-container1:
  docker run -ti -p 8001-8010:8001-8010 kalmhq/echoserver:latest

dp-container2:
  # This container should never receive any request. It should be ensured by abc's TrafficRoute policy
  docker run -p 9941:5678 hashicorp/http-echo -text="ERROR!! I'm a leftover container for a service. I do not have the right koyeb.com/global-deployment tag, hence I should not receive any request!"

dp dp_name="dp" mesh="abc": _build-dp
  @just _inject_ca {{mesh}}

  # Upsert default mesh
  cat koyeb/samples/default-mesh.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -

  # Upsert abc mesh
  cat ./koyeb/samples/{{mesh}}-mesh.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  # Upsert abc Virtual Outbound
  cat ./koyeb/samples/abc-virtual-outbound.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  # Upsert abc TrafficRoute
  cat ./koyeb/samples/abc-traffic-route.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  # Upsert abc MeshTrace
  cat koyeb/samples/abc-meshtrace.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -

  # Wait some time for the mesh to be propagated to the zonal CP...
  sleep 2

  # Finally, upsert the dataplane
  cat ./koyeb/samples/{{mesh}}-{{dp_name}}.yaml | {{kumactl}} apply --config-file koyeb/samples/kumactl-configs/par1-cp.yaml -f -


  # Generate a token for the dataplane instance
  {{kumactl}} generate dataplane-token -m {{mesh}} --valid-for 720h --config-file koyeb/samples/kumactl-configs/par1-cp.yaml > /tmp/{{mesh}}-{{dp_name}}-token
  # Run the dataplane
  {{dev-kuma-dp}} run --dataplane-token-file /tmp/{{mesh}}-{{dp_name}}-token --name {{dp_name}} --log-level info --cp-address https://localhost:5678 --mesh {{mesh}}
