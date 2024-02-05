package handler

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	hook "github.com/no8ge/core/internal/service"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Inject(c *gin.Context) {

	var (
		oldPod          = corev1.Pod{}
		admissionReview = admissionv1.AdmissionReview{}
		patchType       = admissionv1.PatchTypeJSONPatch
	)

	c.BindJSON(&admissionReview)
	json.Unmarshal(admissionReview.Request.Object.Raw, &oldPod)

	protocol := oldPod.ObjectMeta.Annotations["atop.io/protocol"]
	reportPath := oldPod.ObjectMeta.Annotations["atop.io/report-path"]

	newPod := oldPod.DeepCopy()
	if protocol == "s3" {
		sidecar, _ := hook.CreateSidecar(hook.SidecarTypeAtop)
		workerContainer := newPod.Spec.Containers[0]
		workerContainer.VolumeMounts = append(workerContainer.VolumeMounts,
			hook.CreateVolumeMount(sidecar.EmptyDir.Name, reportPath))
		newPod.Spec.Volumes = append(newPod.Spec.Volumes, sidecar.EmptyDir)
		newPod.Spec.Containers = []corev1.Container{sidecar.Container, workerContainer}
	}

	patch, err := hook.DiffPodPatch(oldPod, *newPod)
	if err != nil {
		log.Printf("failed to DiffPodPatch: %v", err)
	}

	admissionReview.Response = &admissionv1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		Result:    &metav1.Status{Status: "Success", Message: "atop.io/sidecar is enable"},
		Patch:     patch,
		PatchType: &patchType,
	}
	c.JSON(200, admissionReview)

}

func Validate(c *gin.Context) {
	req := admissionv1.AdmissionReview{}
	resp := admissionv1.AdmissionReview{
		TypeMeta: req.TypeMeta,
	}
	c.BindJSON(&req)

	pd := corev1.Pod{}
	resp.TypeMeta = req.TypeMeta

	json.Unmarshal(req.Request.Object.Raw, &pd)

	protocol := pd.ObjectMeta.Annotations["atop.io/protocol"]
	reportPath := pd.ObjectMeta.Annotations["atop.io/report-path"]

	if protocol == "" || reportPath == "" {
		resp.Response = &admissionv1.AdmissionResponse{
			UID:     req.Request.UID,
			Allowed: false,
			Result: &metav1.Status{
				Status:  "Failure",
				Message: "if atop.io/enable is true , must be set atop.io/report-path and atop.io/protocol in annotations",
				Reason:  metav1.StatusReason("annotations is not set"),
				Code:    402,
			},
		}
	}
	if protocol != "" && reportPath != "" {
		resp.Response = &admissionv1.AdmissionResponse{
			UID:     req.Request.UID,
			Allowed: true,
			Result: &metav1.Status{
				Status:  "Success",
				Message: "atop.io/enable is true",
				Reason:  metav1.StatusReason("annotations is not set"),
			},
		}
	}
	c.JSON(200, resp)
}
