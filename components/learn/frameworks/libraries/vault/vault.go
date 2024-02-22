package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

func secretReadWrite(client *api.Client) (err error) {
	if secret, err = client.Logical().List("/secret/kv"); err == nil {
		fmt.Println("List:", secret.Data)
	}
	return err
}
