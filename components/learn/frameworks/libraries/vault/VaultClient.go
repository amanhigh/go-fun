package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// HACK: #B Vault Move to Clients
type VaultClientInterface interface {
	GenerateKey(name, algo string, exportable bool) (secret *api.Secret, err error)
	GetKey(name string) (secret *api.Secret, err error)
	ExportKey(name string) (secret *api.Secret, err error)
	RotateKey(name string) (secret *api.Secret, err error)
}

// https://github.com/hashicorp/vault-client-go
// https://github.com/hashicorp/vault-examples
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

/*
*
name - Name of the Key
algo - rsa-4096/2048 - Asymmetric, aes256-gcm96 - Symmetric
*/
func (vc VaultClient) GenerateKey(name, algo string, exportable bool) (secret *api.Secret, err error) {
	keyPath := getKeyPath(name)
	secret, err = vc.client.Logical().Write(keyPath, map[string]any{
		"exportable": exportable,
		"type":       algo,
	})

	if err == nil {
		return vc.GetKey(name)
	}
	return
}

func (vc VaultClient) GetKey(name string) (secret *api.Secret, err error) {
	secret, err = vc.client.Logical().Read(getKeyPath(name))
	return
}

func (vc VaultClient) ExportKey(name string) (secret *api.Secret, err error) {
	keyExportPath := fmt.Sprintf("/transit/export/encryption-key/%v/latest", name)
	secret, err = vc.client.Logical().Read(keyExportPath)
	return
}

func (vc VaultClient) RotateKey(name string) (secret *api.Secret, err error) {
	keyRotatePath := fmt.Sprintf("/transit/keys/%v/rotate", name)
	secret, err = vc.client.Logical().Write(keyRotatePath, nil)
	return
}

func getKeyPath(name string) string {
	return fmt.Sprintf("transit/keys/%v", name)
}
