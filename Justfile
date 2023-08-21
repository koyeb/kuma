
set dotenv-load

artifacts := "./build/artifacts-darwin-amd64/"

default:
  @just --list

run-db:
  -docker run -p 5532:5432 --rm --name kuma-db -e POSTGRES_USER=kuma -e POSTGRES_PASSWORD=kuma -e POSTGRES_DB=kuma -d postgres:14
  sleep 5
  docker run --rm --link kuma-db:kuma-db -e PGHOST=kuma-db -e PGUSER=kuma -e PGPASSWORD=kuma postgres:14 /bin/bash -c "echo 'CREATE DATABASE \"cp-global\" ENCODING UTF8;' | psql"
  docker run --rm --link kuma-db:kuma-db -e PGHOST=kuma-db -e PGUSER=kuma -e PGPASSWORD=kuma postgres:14 /bin/bash -c "echo 'CREATE DATABASE \"cp-par1\" ENCODING UTF8;' | psql"

stop-db:
  docker stop kuma-db


_build-cp:
  make build/kuma-cp

_build-dp:
  make build/kuma-dp
  make build/kumactl

#_gen-root-ca ca_dir="./build/koyeb":
#  @mkdir -p {{ca_dir}}
#  @if [ "{{path_exists(ca_dir + "/rootCA.key")}}" = "false" ] ; then \
#    echo "generating root ca {{ca_dir}}" && \
#    openssl genrsa -out {{ca_dir}}/rootCA.key 4096 && \
#    openssl req -x509 -new -nodes -key {{ca_dir}}/rootCA.key -sha256 -days 10024 -out {{ca_dir}}/rootCA.crt -subj "/C=US/ST=CA/O=koyeb/CN=localhost" ;\
#  fi
#
#_gen-tls dir ca_dir="./build/koyeb" out="cert" config="domains.ext":
#  @mkdir -p {{dir}}
#  @if [ "{{path_exists(dir +"/" + out + "-key.pem")}}" = "false" ] ; then \
#    echo "generating certificated {{dir}}" && \
#    openssl genrsa -out {{dir}}/{{out}}.key 2048 && \
#    openssl req -new -sha256 -key {{dir}}/{{out}}.key -subj "/C=US/ST=CA/O=koyeb/CN=localhost" -out {{dir}}/{{out}}.csr  && \
#    openssl x509 -req -in {{dir}}/{{out}}.csr -CA {{ca_dir}}/rootCA.crt -CAkey {{ca_dir}}/rootCA.key -CAcreateserial -extfile ./koyeb/samples/{{config}} -out {{dir}}/{{out}}.crt -days 5000 -sha256 && \
#    openssl x509 -in {{dir}}/{{out}}.crt -out {{dir}}/{{out}}.pem -text && \
#    openssl x509 -in {{ca_dir}}/rootCA.crt -out {{dir}}/ca.pem -text && \
#    openssl rsa -in {{dir}}/{{out}}.key -text > {{dir}}/{{out}}-key.pem && \
#    openssl x509 -req -in {{dir}}/{{out}}.csr -CA {{ca_dir}}/rootCA.crt -CAkey {{ca_dir}}/rootCA.key -CAcreateserial -extfile ./koyeb/samples/domains.ext -out {{dir}}/{{out}}-client.crt -days 5000 -sha256 && \
#    openssl x509 -in {{dir}}/{{out}}.crt -out {{dir}}/{{out}}-client.pem -text ; \
#  fi
#
#gen-certs: _gen-root-ca (_gen-tls "./build/koyeb/tls-cert") (_gen-tls "./build/koyeb/gateway-cert" "./build/koyeb" "tls" "spiffe.ext")
#  @mkdir -p ./build/koyeb/gateway-client
#  @cp ./build/koyeb/gateway-cert/ca.pem ./build/koyeb/gateway-client/ca.pem
#  @cp ./build/koyeb/gateway-cert/tls-client.pem ./build/koyeb/gateway-client/cert.pem
#  @cp ./build/koyeb/gateway-cert/tls-key.pem ./build/koyeb/gateway-client/cert-key.pem
#  @cp ./build/koyeb/gateway-cert/tls.pem ./build/koyeb/gateway-cert/gateway.pem
#  @cp ./build/koyeb/gateway-cert/tls-key.pem ./build/koyeb/gateway-cert/gateway-key.pem
#
#rm-certs:
#  rm -rf ./build/koyeb/rootCA*
#  rm -rf ./build/koyeb/tls-cert
#  rm -rf ./build/koyeb/gateway-cert
#  rm -rf ./build/koyeb/gateway-client
#
#
#cp-local: gen-certs _build-cp
#  {{artifacts}}/kuma-cp/kuma-cp -c koyeb/samples/config-cp.yaml migrate up
#  TRACING_URL=http://localhost:14268/api/traces \
#  {{artifacts}}/kuma-cp/kuma-cp run --log-level info -c koyeb/samples/config-cp.yaml
#
#
#cp-global: gen-certs _build-cp
#  {{artifacts}}/kuma-cp/kuma-cp -c koyeb/samples/config-global.yaml migrate up
#
#
#  # custom ports for cp-global
#  KUMA_API_SERVER_HTTP_PORT=4281 \
#    KUMA_DIAGNOSTICS_SERVER_PORT=4280 \
#    KUMA_API_SERVER_HTTPS_PORT=4249 \
#    KUMA_DP_SERVER_PORT=4245 \
#    KUMA_XDS_SERVER_GRPC_PORT=4245 \
#    TRACING_URL=http://localhost:14268/api/traces \
#  {{artifacts}}/kuma-cp/kuma-cp run --log-level info -c koyeb/samples/config-global.yaml
#
#
#_init-default:
#  cat ./koyeb/samples/default-mesh.yaml | {{artifacts}}/kumactl/kumactl apply --config-file ./koyeb/samples/kumactl-global-config.yaml -f -

ingress: _build-dp
  {{artifacts}}/kumactl/kumactl generate zone-token --zone=par1 --valid-for 720h --scope ingress > /tmp/dp-token-ingress
  {{artifacts}}/kuma-dp/kuma-dp run \
    --proxy-type=ingress \
    --cp-address=https://127.0.0.1:5678 \
    --dataplane-token-file=/tmp/dp-token-ingress \
    --dataplane-file=koyeb/samples/ingress-par1.yaml


#glb: gen-certs _build-dp _init-default
#  {{artifacts}}/kumactl/kumactl generate zone-ingress-token --zone par1  > /tmp/dp-token-glb
#  {{artifacts}}/kuma-dp/kuma-dp run --dataplane-token-file /tmp/dp-token-glb --dns-coredns-config-template-path ./koyeb/samples/Corefile --dns-coredns-port 10053 --dns-envoy-port 10050 --log-level info --cp-address https://localhost:5678 --ca-cert-file ./build/koyeb/tls-cert/ca.pem --admin-port 4243 -d ./koyeb/samples/ingress-glb.yaml  --proxy-type ingress
#
#igw: gen-certs _build-dp _init-default
#  {{artifacts}}/kumactl/kumactl generate zone-ingress-token --zone par1  > /tmp/dp-token-ingress
#  {{artifacts}}/kuma-dp/kuma-dp run --dataplane-token-file /tmp/dp-token-ingress --dns-coredns-config-template-path ./koyeb/samples/Corefile --dns-coredns-port 10053 --dns-envoy-port 10050 --log-level info --cp-address https://localhost:5678 --ca-cert-file ./build/koyeb/tls-cert/ca.pem --admin-port 4242 -d ./koyeb/samples/ingress.yaml  --proxy-type ingress
#
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
#
dp-container:
  docker run -ti -p 8001-8010:8001-8010 kalmhq/echoserver:latest

dp dp_name="dp" mesh="abc": _build-dp
  # Upsert abc mesh
  cat ./koyeb/samples/{{mesh}}-mesh.yaml | {{artifacts}}/kumactl/kumactl apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  # Upsert abc Virtual Outbound
  cat ./koyeb/samples/abc-virtual-outbound.yaml | {{artifacts}}/kumactl/kumactl apply --config-file koyeb/samples/kumactl-configs/global-cp.yaml -f -
  # Finally, upsert the dataplane
  cat ./koyeb/samples/{{mesh}}-{{dp_name}}.yaml | ./{{artifacts}}/kumactl/kumactl apply --config-file koyeb/samples/kumactl-configs/par1-cp.yaml -f -

  # Generate a token for the dataplane instance
  {{artifacts}}/kumactl/kumactl generate dataplane-token -m {{mesh}} --valid-for 720h --config-file koyeb/samples/kumactl-configs/par1-cp.yaml > /tmp/{{mesh}}-{{dp_name}}-token
  # Run the dataplane
  {{artifacts}}/kuma-dp/kuma-dp run --dataplane-token-file /tmp/{{mesh}}-{{dp_name}}-token --name {{dp_name}} --log-level info --cp-address https://localhost:5678 --mesh {{mesh}}
