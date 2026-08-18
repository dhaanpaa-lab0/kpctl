package platform

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

import (
	"sigs.k8s.io/yaml"
)

// externalSecretGVR identifies the ExternalSecret CRD served by the
// external-secrets.io API group.
var externalSecretGVR = schema.GroupVersionResource{
	Group:    "external-secrets.io",
	Version:  "v1",
	Resource: "externalsecrets",
}

// ExternalSecretParams holds the values that get substituted into the
// ExternalSecret manifest (these correspond to the %%placeholders%% in the
// original template).
type ExternalSecretParams struct {
	ExternalSecretName string // %%externalSecretName%%
	K8sNamespace       string // %%k8sNamespace%%
	SecretStoreKey     string // %%secretStoreKey%%
	SecretStoreName    string // %%secretStoreName%%
}

// validate does simple sanity checking so we fail fast with a useful error
// rather than sending a half-formed manifest to the API server.
func (p ExternalSecretParams) validate() error {
	switch {
	case p.ExternalSecretName == "":
		return fmt.Errorf("externalSecretName is required")
	case p.K8sNamespace == "":
		return fmt.Errorf("k8sNamespace is required")
	case p.SecretStoreKey == "":
		return fmt.Errorf("secretStoreKey is required")
	case p.SecretStoreName == "":
		return fmt.Errorf("secretStoreName is required")
	}
	return nil
}

// buildExternalSecretObject builds the ExternalSecret as an
// unstructured.Unstructured object. Using unstructured (rather than a typed
// external-secrets struct) means this code has no compile-time dependency on
// the external-secrets.io Go module -- it only needs client-go /
// apimachinery, which is convenient since ExternalSecret is a CRD, not a
// built-in Kubernetes type.
func buildExternalSecretObject(p ExternalSecretParams) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "external-secrets.io/v1",
			"kind":       "ExternalSecret",
			"metadata": map[string]any{
				"name":      p.ExternalSecretName,
				"namespace": p.K8sNamespace,
				"finalizers": []any{
					"externalsecrets.external-secrets.io/externalsecret-cleanup",
				},
			},
			"spec": map[string]any{
				"refreshInterval": "15s",
				"dataFrom": []any{
					map[string]any{
						"extract": map[string]any{
							"conversionStrategy": "Default",
							"decodingStrategy":   "None",
							"key":                p.SecretStoreKey,
							"metadataPolicy":     "None",
						},
					},
				},
				"secretStoreRef": map[string]any{
					"kind": "ClusterSecretStore",
					"name": p.SecretStoreName,
				},
				"target": map[string]any{
					"name":           p.ExternalSecretName,
					"creationPolicy": "Owner",
					"deletionPolicy": "Retain",
				},
			},
		},
	}
	return obj
}

// GenerateExternalSecretYAML renders the ExternalSecret manifest as YAML text,
// substituting the given parameters for the %%placeholders%% in the original
// template. The returned string is suitable for piping into
// `kubectl apply -f -`, e.g.:
//
//	yamlText, err := GenerateExternalSecretYAML(params)
//	cmd := exec.Command("kubectl", "apply", "-f", "-")
//	cmd.Stdin = strings.NewReader(yamlText)
//	err = cmd.Run()
func GenerateExternalSecretYAML(p ExternalSecretParams) (string, error) {
	if err := p.validate(); err != nil {
		return "", fmt.Errorf("invalid ExternalSecret params: %w", err)
	}

	obj := buildExternalSecretObject(p)

	out, err := yaml.Marshal(obj.Object)
	if err != nil {
		return "", fmt.Errorf("marshal ExternalSecret to yaml: %w", err)
	}

	return string(out), nil
}

// ApplyExternalSecret applies the ExternalSecret directly to the cluster
// using the Kubernetes dynamic client (server-side apply), which is
// necessary because ExternalSecret is a CRD and not a built-in type known to
// client-go's typed clientset.
//
// dynClient should be constructed by the caller, e.g.:
//
//	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
//	dynClient, err := dynamic.NewForConfig(cfg)
//
// fieldManager identifies the owner of the applied fields for server-side
// apply purposes (e.g. "my-controller"); pass a stable, non-empty value.
func ApplyExternalSecret(ctx context.Context, dynClient dynamic.Interface, p ExternalSecretParams, fieldManager string) (*unstructured.Unstructured, error) {
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("invalid ExternalSecret params: %w", err)
	}
	if dynClient == nil {
		return nil, fmt.Errorf("dynClient must not be nil")
	}
	if fieldManager == "" {
		fieldManager = "externalsecret-applier"
	}

	obj := buildExternalSecretObject(p)

	data, err := obj.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal ExternalSecret to json: %w", err)
	}

	applied, err := dynClient.
		Resource(externalSecretGVR).
		Namespace(p.K8sNamespace).
		Patch(
			ctx,
			p.ExternalSecretName,
			// Server-side apply: creates the object if it doesn't exist,
			// or updates the fields we own if it does.
			types.ApplyPatchType,
			data,
			metav1.PatchOptions{
				FieldManager: fieldManager,
				Force:        new(true),
			},
		)
	if err != nil {
		return nil, fmt.Errorf("apply ExternalSecret %s/%s: %w", p.K8sNamespace, p.ExternalSecretName, err)
	}

	return applied, nil
}
