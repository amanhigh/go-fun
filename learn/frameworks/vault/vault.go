package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

func VaultFun() {
	var data = map[string]interface{}{
		"Id":       1,
		"Name":     "Aman",
		"password": "Preet",
	}

	var err error

	var client *api.Client
	if client, err = api.NewClient(&api.Config{Address: "http://docker:8200"}); err == nil {
		client.SetToken("root-token")
		err = secretReadWrite(client, data)

	}

	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func secretReadWrite(client *api.Client, data map[string]interface{}) (err error) {
	var secret *api.Secret
	path := "/secret/kv/test"
	if secret, err = client.Logical().Write(path, data); err == nil {
		fmt.Println("Write Complete", secret)
	}
	if secret, err = client.Logical().Read(path); err == nil {
		fmt.Println("Read ", secret.Data, secret.LeaseDuration)
	}

	if secret, err = client.Logical().List("/secret/kv"); err == nil {
		fmt.Println("List:", secret.Data)
	}
	return err
}
