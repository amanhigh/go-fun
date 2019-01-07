package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type VaultClientInterface interface {
	GenerateKey(name, algo string, exportable bool) (secret *api.Secret, err error)
	GetKey(name string) (secret *api.Secret, err error)
}

type VaultClient struct {
	client *api.Client
}

func NewVaultClient(host, token string, port int) (vaultClient VaultClientInterface, err error) {
	vaultUrl := fmt.Sprintf("http://%v:%v", host, port)

	var client *api.Client
	if client, err = api.NewClient(&api.Config{Address: vaultUrl}); err == nil {
		client.SetToken(token)
		vaultClient = VaultClient{client}
	}
	return
}

/**
name - Name of the Key
algo - rsa-4096/2048 - Asymmetric, aes256-gcm96 - Symmetric
*/
func (self VaultClient) GenerateKey(name, algo string, exportable bool) (secret *api.Secret, err error) {
	keyPath := getKeyPath(name)
	secret, err = self.client.Logical().Write(keyPath, map[string]interface{}{
		"exportable": exportable,
		"type":       algo,
	})

	if err == nil {
		return self.GetKey(name)
	}
	return
}

func (self VaultClient) GetKey(name string) (secret *api.Secret, err error) {
	secret, err = self.client.Logical().Read(getKeyPath(name))
	return
}

func getKeyPath(name string) string {
	return fmt.Sprintf("transit/keys/%v", name)
}
