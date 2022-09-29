package kind

import (
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"strings"
)

const (
	VMKind    = "VirtualMachine"
	HostLable = "host"
	// todo
	HostLableDefaultValue = "vm"
)

var (
	vmRequiredLabels = []string{
		HostLable,
	}
	vmAddLabels = map[string]string{
		HostLable: HostLableDefaultValue,
	}
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	newValues := make(map[string]string)
	updateValues := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			newValues[key] = value
		} else if target[key] != added[key] {
			updateValues[key] = value
		}
	}
	klog.Infof("update label:%+v", updateValues)
	for k, v := range newValues {
		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  "/metadata/labels/" + strings.ReplaceAll(k, "/", "~1"),
			Value: v,
		})
	}
	// fix patch key with /: use ~1 to encode /
	for k, v := range updateValues {
		patch = append(patch, patchOperation{
			Op:    "replace",
			Path:  "/metadata/labels/" + strings.ReplaceAll(k, "/", "~1"),
			Value: v,
		})
	}

	return patch
}

func createPatch(availableLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation
	klog.Infof("availableLabels: %+v, labels:%+v", availableLabels, labels)
	patch = append(patch, updateLabels(availableLabels, labels)...)

	return json.Marshal(patch)
}

func Addlabel(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var availableLabels map[string]string

	klog.Infof("======begin Mutating Admission for Namespace=[%v], Kind=[%v], Name=[%v]======", req.Namespace, req.Kind.Kind, req.Name)

	switch req.Kind.Kind {
	case VMKind:
		klog.Infof("vmkind str:%s: %+v", string(req.Object.Raw))
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			klog.Infof("Could not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		klog.Infof("deployment fake:%+v", deployment)
		availableLabels = deployment.Labels
	default:
		msg := fmt.Sprintf("Not support for this Kind of resource  %v", req.Kind.Kind)
		klog.Info(msg)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}

	var patchBytes []byte
	var err error
	// add labels and annotation
	if req.Kind.Kind == VMKind {
		patchBytes, err = createPatch(availableLabels, vmAddLabels)
		if err != nil {
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
	}

	klog.Infof("AdmissionResponse: patch=%v", string(patchBytes))
	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
