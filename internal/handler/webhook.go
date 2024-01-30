package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/wI2L/jsondiff"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Inject(c *gin.Context) {

	req := admissionv1.AdmissionReview{}
	resp := admissionv1.AdmissionReview{
		TypeMeta: req.TypeMeta,
	}
	c.BindJSON(&req)

	pod := corev1.Pod{}
	sidecarContainer := corev1.Container{}
	json.Unmarshal(req.Request.Object.Raw, &pod)

	reportPath := pod.ObjectMeta.Annotations["atop.io/report-path"]

	newPod := pod.DeepCopy()
	workerContainer := newPod.Spec.Containers[0]
	workerContainer.VolumeMounts = append(workerContainer.VolumeMounts, corev1.VolumeMount{
		Name:      "cache-volume",
		MountPath: reportPath,
	})
	newPod.Spec.Containers[0] = workerContainer
	protocol := pod.ObjectMeta.Annotations["atop.io/protocol"]
	if protocol == "s3" {
		sidecarContainer.Name = "sidecar"
		sidecarContainer.Image = "no8ge/sidecar:1.0.0"
		sidecarContainer.ImagePullPolicy = "Always"
		sidecarContainer.Command = append(sidecarContainer.Command, "/bin/sh")
		sidecarContainer.Args = append(sidecarContainer.Args, "-c")
		sidecarContainer.Args = append(sidecarContainer.Args, "mc alias set atop http://$MINIO_HOST $MINIO_ACCESS_KEY $MINIO_SECRET_KEY; mc mirror --remove --watch --overwrite /data atop/result/$REPORT")
		sidecarContainer.Env = append(sidecarContainer.Env, corev1.EnvVar{
			Name: "REPORT",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		})
		sidecarContainer.Env = append(sidecarContainer.Env, corev1.EnvVar{
			Name:  "MINIO_ACCESS_KEY",
			Value: "admin",
		})
		sidecarContainer.Env = append(sidecarContainer.Env, corev1.EnvVar{
			Name:  "MINIO_SECRET_KEY",
			Value: "changeme",
		})
		sidecarContainer.Env = append(sidecarContainer.Env, corev1.EnvVar{
			Name:  "MINIO_HOST",
			Value: "files-minio.default:9000",
		})
		sidecarContainer.VolumeMounts = append(sidecarContainer.VolumeMounts, corev1.VolumeMount{
			Name:      "cache-volume",
			MountPath: "/data",
		})
	}

	newPod.Spec.Volumes = append(newPod.Spec.Volumes, corev1.Volume{
		Name: "cache-volume",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})
	newPod.Spec.Containers = append(newPod.Spec.Containers, sidecarContainer)

	diff, err := jsondiff.Compare(pod, newPod)
	if err != nil {
		return
	}

	patch, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		return
	}

	patchType := admissionv1.PatchTypeJSONPatch
	resp.Response = &admissionv1.AdmissionResponse{
		UID:       req.Request.UID,
		Allowed:   true,
		Result:    &metav1.Status{Status: "Success", Message: "atop.io/sidecar is enable"},
		Patch:     patch,
		PatchType: &patchType,
	}

	c.JSON(200, resp)

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
