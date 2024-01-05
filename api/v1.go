package api

import (
	"encoding/json"
	"log"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"
)

func V1(r *gin.Engine) {

	plugins := r.Group("/v1")
	{
		plugins.GET("/", func(c *gin.Context) {
			c.String(200, "Hello, Geektutu")
		})
		plugins.GET("/info", func(c *gin.Context) {
			c.String(200, "Hello, Geektutu")
		})
	}

	webHook := r.Group("/v1")
	{
		webHook.POST("/validate", func(c *gin.Context) {
			req := admissionv1.AdmissionReview{}
			c.BindJSON(&req)

			pd := corev1.Pod{}
			resp := admissionv1.AdmissionReview{
				TypeMeta: req.TypeMeta,
			}

			json.Unmarshal(req.Request.Object.Raw, &pd)

			protocol := pd.ObjectMeta.Annotations["atop.io/protocol"]
			reportPath := pd.ObjectMeta.Annotations["atop.io/report-path"]

			if protocol == "" || reportPath == "" {
				log.Println("atop.io/sidecar value is not valid, it must be enable")
				resp.Response = &admissionv1.AdmissionResponse{
					UID:     req.Request.UID,
					Allowed: false,
					Result: &metav1.Status{
						Status:  "Failure",
						Message: "if sidecar is enable , must be set atop.io/report-path and atop.io/protocol in annotations",
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
						Message: "atop.io/sidecar is enable",
					},
				}
			}
			log.Println(resp)
			c.JSON(200, resp)
		})
		webHook.POST("/inject", func(c *gin.Context) {
			c.String(200, "Hello, Geektutu")
		})
	}

}
