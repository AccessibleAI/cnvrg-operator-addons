package test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	deployValues = map[string]string{"clusterDomain": "aws.dilerous.cloud",
		"controlPlane.image":      "cnvrg/app:v4.7.90",
		"registry.user":           "cnvrghelm",
		"registry.password":       "abbb7835-fef8-42af-be0f-ef6750bde5a0",
		"networking.ingress.type": "ingress"}
)

func TestDeployCnvrg(t *testing.T) {
	// Path to the helm chart we will test
	//helmChartPath := "../charts/minimal-pod"

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// We also specify that we are working in the default namespace (required to get the Pod)
	kubectlOptions := k8s.NewKubectlOptions("", "", "cnvrg")

	// Setup the args. For this test, we will set the following input values:
	// - image=nginx:1.15.8
	options := &helm.Options{
		SetValues:      deployValues,
		KubectlOptions: k8s.NewKubectlOptions("", "", "cnvrg"),
	}

	// We generate a unique release name so that we can refer to after deployment.
	// By doing so, we can schedule the delete call here so that at the end of the test, we run
	// `helm delete RELEASE_NAME` to clean up any resources that were created.
	releaseName := fmt.Sprintf("cnvrg-%s", strings.ToLower(random.UniqueId()))
	defer helm.Delete(t, options, releaseName, true)

	// Deploy the chart using `helm install`. Note that we use the version without `E`, since we want to assert the
	// install succeeds without any errors.
	helm.Install(t, options, DeployPath, releaseName)

	// Now that the chart is deployed, verify the deployment. This function will open a tunnel to the Pod and hit the
	// nginx container endpoint.
	//podName := fmt.Sprintf("%s-minimal-pod", releaseName)
	podName := getDeployments()
	fmt.Println(podName)
	verifyAppPod(t, kubectlOptions, podName)
}

func getDeployments() string {

	var rules = clientcmd.NewDefaultClientConfigLoadingRules()
	var kubeconfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Error creating Kubernetes config: %v", err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	// List all deployments in the cluster
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing deployments: %v", err)
	}
	for _, deployment := range deployments.Items {
		fmt.Printf("Found deployment %s in namespace %s\n", deployment.Name, deployment.Namespace)
	}

	deployment, err := clientset.AppsV1().Deployments("cnvrg").Get(context.Background(), "app", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error getting deployment: %v", err)
	}

	// Get the selector labels of the deployment
	selector := deployment.Spec.Selector
	labelSelector, err := v1.LabelSelectorAsSelector(selector)
	if err != nil {
		log.Fatalf("Error getting label selector: %v", err)
	}

	// Use the label selector to find the pod
	podList, err := clientset.CoreV1().Pods("cnvrg").List(context.Background(), v1.ListOptions{
		LabelSelector: labelSelector.String(),
	})
	if err != nil {
		log.Fatalf("Error getting pod list: %v", err)
	}

	if len(podList.Items) > 0 {
		// Print the name of the first pod in the list
		fmt.Printf("Pod name: %s\n", podList.Items[0].Name)
	} else {
		log.Fatalf("No pods found for deployment %s", deployment.Name)
	}

	return podList.Items[0].Name

}

func verifyAppPod(t *testing.T, kubectlOptions *k8s.KubectlOptions, podName string) {
	// Wait for the pod to come up. It takes some time for the Pod to start, so retry a few times.
	retries := 15
	sleep := 5 * time.Second
	k8s.WaitUntilPodAvailable(t, kubectlOptions, podName, retries, sleep)

	// We will first open a tunnel to the pod, making sure to close it at the end of the test.
	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypePod, podName, 0, 8080)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// ... and now that we have the tunnel, we will verify that we get back a 200 OK with the nginx welcome page.
	// It takes some time for the Pod to start, so retry a few times.
	endpoint := fmt.Sprintf("http://%s", tunnel.Endpoint())
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		endpoint,
		nil,
		retries,
		sleep,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)
}

//&& strings.Contains(body, "Welcome to nginx")
