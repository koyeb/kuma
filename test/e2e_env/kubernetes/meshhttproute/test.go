package meshhttproute

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/plugins/policies/meshhttproute/api/v1alpha1"
	. "github.com/kumahq/kuma/test/framework"
	"github.com/kumahq/kuma/test/framework/client"
	"github.com/kumahq/kuma/test/framework/deployments/testserver"
	"github.com/kumahq/kuma/test/framework/envs/kubernetes"
)

func Test() {
	meshName := "meshhttproute"
	namespace := "meshhttproute"

	BeforeAll(func() {
		err := NewClusterSetup().
			Install(MeshKubernetes(meshName)).
			Install(NamespaceWithSidecarInjection(namespace)).
			Install(testserver.Install(
				testserver.WithName("test-client"),
				testserver.WithMesh(meshName),
				testserver.WithNamespace(namespace),
			)).
			Install(testserver.Install(
				testserver.WithName("test-server"),
				testserver.WithMesh(meshName),
				testserver.WithNamespace(namespace),
			)).
			Install(testserver.Install(
				testserver.WithNamespace(namespace),
				testserver.WithName("external-service"),
			)).
			Setup(kubernetes.Cluster)
		Expect(err).ToNot(HaveOccurred())

		Expect(
			k8s.RunKubectlE(kubernetes.Cluster.GetTesting(), kubernetes.Cluster.GetKubectlOptions(), "delete", "trafficroute", "route-all-meshhttproute"),
		).To(Succeed())
	})

	E2EAfterEach(func() {
		Expect(DeleteMeshResources(kubernetes.Cluster, meshName, v1alpha1.MeshHTTPRouteResourceTypeDescriptor)).To(Succeed())
		Expect(DeleteMeshResources(kubernetes.Cluster, meshName, core_mesh.ExternalServiceResourceTypeDescriptor)).To(Succeed())
	})

	E2EAfterAll(func() {
		Expect(kubernetes.Cluster.TriggerDeleteNamespace(namespace)).To(Succeed())
		Expect(kubernetes.Cluster.DeleteMesh(meshName)).To(Succeed())
	})

	It("should use MeshHTTPRoute if any MeshHTTPRoutes are present", func() {
		Eventually(func(g Gomega) {
			_, err := client.CollectResponse(kubernetes.Cluster, "test-client", "test-server_meshhttproute_svc_80.mesh", client.FromKubernetesPod(namespace, "test-client"))
			g.Expect(err).To(HaveOccurred())
		}, "30s", "1s").Should(Succeed())

		// when
		Expect(YamlK8s(fmt.Sprintf(`
apiVersion: kuma.io/v1alpha1
kind: MeshHTTPRoute
metadata:
  name: route-1
  namespace: %s
  labels:
    kuma.io/mesh: %s
spec:
  targetRef:
    kind: MeshService
    name: test-client_%s_svc_80
  to:
    - targetRef:
        kind: MeshService
        name: nonexistent-service-that-activates-default
      rules: []
`, Config.KumaNamespace, meshName, meshName))(kubernetes.Cluster)).To(Succeed())

		Eventually(func(g Gomega) {
			response, err := client.CollectResponse(kubernetes.Cluster, "test-client", "test-server_meshhttproute_svc_80.mesh", client.FromKubernetesPod(namespace, "test-client"))
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(response.Instance).To(HavePrefix("test-server"))
		}, "30s", "1s").Should(Succeed())
	})

	It("should split traffic between internal and external services", func() {
		// given
		Expect(kubernetes.Cluster.Install(YamlK8s(fmt.Sprintf(`
apiVersion: kuma.io/v1alpha1
kind: ExternalService
metadata:
  name: external-service-1
  namespace: %s
mesh: %s
spec:
  tags:
    kuma.io/service: external-service
    kuma.io/protocol: http
  networking:
    address: external-service.%s.svc.cluster.local:80 # .svc.cluster.local is needed, otherwise Kubernetes will resolve this to the real IP
`, Config.KumaNamespace, meshName, namespace)))).To(Succeed())

		// when
		Expect(YamlK8s(fmt.Sprintf(`
apiVersion: kuma.io/v1alpha1
kind: MeshHTTPRoute
metadata:
  name: route-2
  namespace: %s
  labels:
    kuma.io/mesh: %s
spec:
  targetRef:
    kind: MeshService
    name: test-client_%s_svc_80
  to:
    - targetRef:
        kind: MeshService
        name: test-server_meshhttproute_svc_80
      rules: 
        - matches:
            - path: 
                type: Prefix
                value: /
          default:
            backendRefs:
              - kind: MeshService
                name: test-server_meshhttproute_svc_80
                weight: 50
              - kind: MeshService
                name: external-service
                weight: 50
`, Config.KumaNamespace, meshName, meshName))(kubernetes.Cluster)).To(Succeed())

		// then receive responses from 'test-server_meshhttproute_svc_80'
		Eventually(func(g Gomega) {
			response, err := client.CollectResponse(kubernetes.Cluster, "test-client", "test-server_meshhttproute_svc_80.mesh", client.FromKubernetesPod(namespace, "test-client"))
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(response.Instance).To(HavePrefix("test-server"))
		}, "30s", "1s").Should(Succeed())

		// and then receive responses from 'external-service'
		Eventually(func(g Gomega) {
			response, err := client.CollectResponse(kubernetes.Cluster, "test-client", "test-server_meshhttproute_svc_80.mesh", client.FromKubernetesPod(namespace, "test-client"))
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(response.Instance).To(HavePrefix("external-service"))
		}, "30s", "1s").Should(Succeed())
	})
}