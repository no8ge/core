package api

import (
	"encoding/json"
	"log"
	"net/http"

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
			var pd corev1.Pod
			var resp admissionv1.AdmissionReview
			resp.TypeMeta = req.TypeMeta

			json.Unmarshal(req.Request.Object.Raw, &pd)

			label, ok := pd.ObjectMeta.Labels["atop.io/sidecar"]
			log.Println(ok)
			if !ok {
				log.Println("ignore sidecar validate")
				resp.Response = &admissionv1.AdmissionResponse{
					UID:     req.Request.UID,
					Allowed: true,
					Result:  nil,
				}
			} else {
				if label != "enable" {
					log.Println("atop.io/sidecar value is not valid, it must be enable")
					resp.Response = &admissionv1.AdmissionResponse{
						UID:     req.Request.UID,
						Allowed: false,
						Result: &metav1.Status{
							Status:  "Failure",
							Message: "if want enable sidecar, key must be atop.io, value must be enable",
							Reason:  metav1.StatusReason("atop.io/sidecar key or value is not valid"),
							Code:    402,
						},
					}
				}
				if label == "enable" {
					log.Println("sidecar validate pass")
					resp.Response = &admissionv1.AdmissionResponse{
						UID:     req.Request.UID,
						Allowed: true,
						Result:  nil,
					}
				}
			}
			log.Println(resp)
			c.JSON(http.StatusOK, resp)
		})
		webHook.POST("/inject", func(c *gin.Context) {
			c.String(200, "Hello, Geektutu")
		})
	}

}
