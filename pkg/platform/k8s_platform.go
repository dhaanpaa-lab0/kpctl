package platform

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetK8SPlatformConfig() *api.Config {
	// Get the default kubeconfig file path
	rules := clientcmd.NewDefaultClientConfigLoadingRules()

	// Read the kubeconfig data
	config, err := rules.Load()
	if err != nil {
		panic(err)
	}

	return config
}
func GetAvailableContexts() []string {
	contextNames := make([]string, 0)
	for name := range GetK8SPlatformConfig().Contexts {
		contextNames = append(contextNames, name)
	}
	return contextNames
}
func GetK8SPlatformRestConfig() *rest.Config {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		panic(err)
	}
	return config
}
func GetNamespaces() []string {
	// 2. Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(GetK8SPlatformRestConfig())
	if err != nil {
		panic(err)
	}

	// 3. List all namespaces using the current context
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	namespaceNames := make([]string, 0)

	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}
	return namespaceNames
}
