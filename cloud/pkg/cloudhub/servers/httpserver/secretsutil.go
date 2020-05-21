package httpserver

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/kubeedge/kubeedge/cloud/pkg/edgecontroller/utils"
)

const (
	NamespaceSystem string = "kubeedge"

	TokenSecretName      string = "tokensecret"
	TokenDataName        string = "tokendata"
	CaSecretName         string = "casecret"
	CloudCoreSecretName  string = "cloudcoresecret"
	CaDataName           string = "cadata"
	CaKeyDataName        string = "cakeydata"
	CloudCoreCertName    string = "cloudcoredata"
	CloudCoreKeyDataName string = "cloudcorekeydata"
)

func GetSecret(secretName string, ns string) (*v1.Secret, error) {
	cli, err := utils.KubeClient()
	if err != nil {
		fmt.Printf("%v", err)
	}
	return cli.CoreV1().Secrets(ns).Get(secretName, metav1.GetOptions{})
}

// CreateSecret creates a secret
func CreateSecret(secret *v1.Secret, ns string) error {
	cli, err := utils.KubeClient()
	if err != nil {
		fmt.Printf("%v", err)
	}
	if _, err := cli.CoreV1().Secrets(ns).Create(secret); err != nil {
		if apierrors.IsAlreadyExists(err) {
			cli.CoreV1().Secrets(ns).Update(secret)
		} else {
			klog.Errorf("Failed to create the secret, namespace: %s, name: %s, err: %v", ns, secret.Name, err)
			return fmt.Errorf("Failed to create the secret, namespace: %s, name: %s, err: %v", ns, secret.Name, err)
		}
	}
	return nil
}

func CreateTokenSecret(caHashAndToken []byte) error {
	token := &v1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenSecretName,
			Namespace: NamespaceSystem,
		},
		Data: map[string][]byte{
			TokenDataName: caHashAndToken,
		},
		StringData: map[string]string{},
		Type:       "Opaque",
	}
	return CreateSecret(token, NamespaceSystem)
}

func CreateCaSecret(certDER, key []byte) error {
	caSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      CaSecretName,
			Namespace: NamespaceSystem,
		},
		Data: map[string][]byte{
			CaDataName:    certDER,
			CaKeyDataName: key,
		},
		StringData: map[string]string{},
		Type:       "Opaque",
	}
	return CreateSecret(caSecret, NamespaceSystem)
}

func CreateCloudCoreSecret(certDER, key []byte) error {
	cloudCoreCert := &v1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      CloudCoreSecretName,
			Namespace: NamespaceSystem,
		},
		Data: map[string][]byte{
			CloudCoreCertName:    certDER,
			CloudCoreKeyDataName: key,
		},
		StringData: map[string]string{},
		Type:       "Opaque",
	}
	return CreateSecret(cloudCoreCert, NamespaceSystem)
}