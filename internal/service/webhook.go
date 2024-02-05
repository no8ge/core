package service

import (
	"encoding/json"
	"log"

	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

type SidecarType string

const (
	SidecarTypeAtop     SidecarType = "atop"
	SidecarTypeFilebeat SidecarType = "filebeat"
)

type Sidecar struct {
	Type             SidecarType
	Container        corev1.Container
	EmptyDir         corev1.Volume
	PersistentVolume corev1.Volume
}

func CreateSidecar(t SidecarType) (*Sidecar, error) {
	restartPolicy := corev1.ContainerRestartPolicyAlways
	envfromMinio := corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "files-config",
			},
		},
	}
	cacheVolumeMount := corev1.VolumeMount{
		Name:      "cache-volume",
		MountPath: "/data",
	}

	c := corev1.Container{
		Name:          "sidecar",
		RestartPolicy: &restartPolicy,
	}

	var (
		sidecar     = &Sidecar{}
		port        corev1.ContainerPort
		cachevolume = corev1.Volume{
			Name: "cache-volume",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		datavolume = corev1.Volume{
			Name: "date-volume",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "data-volume",
				},
			},
		}
	)

	switch t {
	case SidecarTypeAtop:
		port = corev1.ContainerPort{
			ContainerPort: 5045,
		}
		c.Image = "no8ge/sidecar:1.0.0"
		c.Command = []string{"/bin/sh"}
		c.Args = []string{"-c", "mc alias set atop http://$MINIO_HOST $MINIO_ACCESS_KEY $MINIO_SECRET_KEY;mc mirror --remove --watch --overwrite /data atop/result/$HOSTNAME"}
		c.Ports = []corev1.ContainerPort{port}
		c.EnvFrom = []corev1.EnvFromSource{envfromMinio}
		c.VolumeMounts = []corev1.VolumeMount{cacheVolumeMount}
		sidecar.Type = t
		sidecar.Container = c
		sidecar.EmptyDir = cachevolume
		sidecar.PersistentVolume = datavolume

	case SidecarTypeFilebeat:
		port = corev1.ContainerPort{
			ContainerPort: 5044,
		}
		c.Image = "no8ge/filebeat:1.0.0"
		c.Ports = []corev1.ContainerPort{port}
		sidecar.Type = t
		sidecar.Container = c
		sidecar.EmptyDir = cachevolume
		sidecar.PersistentVolume = datavolume
	}
	return sidecar, nil
}

func CreateVolumeMount(name string, path string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		MountPath: path,
	}
}

func DiffPodPatch(oldPod, newPod corev1.Pod) ([]byte, error) {
	diff, err := jsondiff.Compare(oldPod, newPod)
	if err != nil {
		log.Printf("failed to compare pod: %v", err)
	}
	patch, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		log.Printf("failed to MarshalIndent: %v", err)
	}
	return patch, err
}
