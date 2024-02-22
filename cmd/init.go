package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/no8ge/core/pkg/k8s"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/admissionregistration/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	days     = 365
	certFile = "./cert/core.pem"
	keyFile  = "./cert/core.key"
	hookName = "core.atop.io"
)

func init() {
	rootCmd.AddCommand(InitCmd)
	InitCmd.PersistentFlags().StringP("namespace", "N", "default", "the namespace of core service in k8s")
	InitCmd.PersistentFlags().StringP("name", "n", "core", "the name of core service in k8s")
}

func generateSelfSignedCert(webhookURL string) error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(days) * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: webhookURL,
		},
		DNSNames:              []string{webhookURL},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	err = os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}), 0644)
	if err != nil {
		return fmt.Errorf("failed to write cert file: %v", err)
	}

	err = os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)}), 0644)
	if err != nil {
		return fmt.Errorf("failed to write key file: %v", err)
	}
	log.Printf("successed to generate self signed cert")
	return nil
}

func createValidatingWebhookConfig(name string, kubeClient kubernetes.Interface, wcc v1.WebhookClientConfig) error {
	var (
		scope       = v1.AllScopes
		sideEffects = v1.SideEffectClassNone
	)

	validatingWebhook := &v1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Webhooks: []v1.ValidatingWebhook{{
			Name:         hookName,
			ClientConfig: wcc,
			Rules: []v1.RuleWithOperations{{
				Operations: []v1.OperationType{
					v1.Create,
					v1.Update,
				},
				Rule: v1.Rule{
					APIGroups:   []string{"*"},
					APIVersions: []string{"*"},
					Resources:   []string{"pods", "deployments"},
					Scope:       &scope,
				},
			}},
			FailurePolicy:     &[]v1.FailurePolicyType{"Fail"}[0],
			NamespaceSelector: &metav1.LabelSelector{},
			ObjectSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      "atop.io/enable",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"true"},
				}},
			},
			SideEffects:             &sideEffects,
			AdmissionReviewVersions: []string{"v1beta1", "v1"},
			TimeoutSeconds:          &[]int32{30}[0],
		}},
	}

	_, err := kubeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.TODO(), validatingWebhook.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		log.Printf("failed to get ValidatingWebhookConfiguration: %v", err)
	} else if err == nil {
		err := kubeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(context.TODO(), validatingWebhook.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("failed to delete ValidatingWebhookConfiguration: %v", err)
			return nil
		}
		log.Printf("successed to delete old ValidatingWebhookConfiguration")
	} else {
		log.Printf("failed to get ValidatingWebhookConfiguration: %v", err)
	}
	_, createErr := kubeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(context.TODO(), validatingWebhook, metav1.CreateOptions{})
	if createErr != nil {
		log.Printf("failed to create ValidatingWebhookConfiguration: %v", createErr)
		return nil
	}
	log.Printf("successed to create ValidatingWebhookConfiguration")
	return nil
}

func createMutatingWebhookConfig(name string, kubeClient kubernetes.Interface, wcc v1.WebhookClientConfig) error {
	var (
		scope              = v1.AllScopes
		sideEffects        = v1.SideEffectClassNone
		reinvocationPolicy = v1.NeverReinvocationPolicy
	)

	mutatingWebhook := &v1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Webhooks: []v1.MutatingWebhook{{
			Name:         hookName,
			ClientConfig: wcc,
			Rules: []v1.RuleWithOperations{{
				Operations: []v1.OperationType{
					v1.Create,
					v1.Update,
				},
				Rule: v1.Rule{
					APIGroups:   []string{"*"},
					APIVersions: []string{"*"},
					Resources:   []string{"pods", "deployments"},
					Scope:       &scope,
				},
			}},
			FailurePolicy:     &[]v1.FailurePolicyType{"Ignore"}[0],
			NamespaceSelector: &metav1.LabelSelector{},
			ObjectSelector: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{{
				Key:      "atop.io/enable",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"true"},
			}}},
			ReinvocationPolicy:      &reinvocationPolicy,
			SideEffects:             &sideEffects,
			AdmissionReviewVersions: []string{"v1beta1", "v1"},
			TimeoutSeconds:          &[]int32{30}[0],
		}},
	}

	_, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.TODO(), mutatingWebhook.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		log.Printf("failed to get MutatingWebhookConfiguration: %v", err)
	} else if err == nil {
		err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.TODO(), mutatingWebhook.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("failed to delete MutatingWebhookConfiguration: %v", err)
			return nil
		}
		log.Printf("successed to delete old MutatingWebhookConfiguration")
	} else {
		log.Printf("failed to get MutatingWebhookConfiguration: %v", err)
	}
	_, createErr := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), mutatingWebhook, metav1.CreateOptions{})
	if createErr != nil {
		log.Printf("failed to create MutatingWebhookConfiguration: %v", createErr)
		return nil
	}
	log.Printf("successed to create MutatingWebhookConfiguration")
	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init core to k8s",
	Long:  `Init core to k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		mod, _ := cmd.Flags().GetString("mod")
		log.Printf("successed to get mod: %v", mod)

		port, _ := cmd.Flags().GetString("port")
		name, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")

		intPort, err := strconv.Atoi(port)
		int32Value := int32(intPort)
		if err != nil {
			panic(err)
		}

		client, err := k8s.Client()
		if err != nil {
			panic(err.Error())
		}

		var (
			webhookHost string
			vc          v1.WebhookClientConfig
			mc          v1.WebhookClientConfig
			vcPath      = "/v1/validate"
			mcPath      = "/v1/inject"
		)
		if mod == "debug" {
			webhookHost = "host.docker.internal"
			vcWebhookURL := fmt.Sprintf("https://%s:%s%s", webhookHost, port, vcPath)
			mcWebhookURL := fmt.Sprintf("https://%s:%s%s", webhookHost, port, mcPath)
			vc.URL = &vcWebhookURL
			mc.URL = &mcWebhookURL
		}
		if mod == "release" {
			webhookHost = fmt.Sprintf("%s.%s.svc", name, namespace)
			vc.Service = &v1.ServiceReference{
				Namespace: namespace,
				Name:      name,
				Path:      &vcPath,
				Port:      &int32Value,
			}
			mc.Service = &v1.ServiceReference{
				Namespace: namespace,
				Name:      name,
				Path:      &mcPath,
				Port:      &int32Value,
			}
		}
		generateSelfSignedCert(webhookHost)
		pemData, err := os.ReadFile(certFile)
		if err != nil {
			panic(err)
		}
		vc.CABundle = pemData
		mc.CABundle = pemData
		createValidatingWebhookConfig(name, client, vc)
		createMutatingWebhookConfig(name, client, mc)
	},
}
