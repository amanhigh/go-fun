package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/dhutil"
)

func VaultFun() {
	var err error

	var client *api.Client
	if client, err = api.NewClient(&api.Config{Address: "http://docker:8200"}); err == nil {
		client.SetToken("root-token")
		//err = secretReadWrite(client)
		err = transitFun(client)

	}

	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func transitFun(client *api.Client) (err error) {
	//	Transit
	var secret *api.Secret
	//Create Key
	secret, err = client.Logical().Write("transit/keys/aman", map[string]interface{}{
		"exportable": true,
		//rsa-4096 - Asymmetric, aes256-gcm96 - Symmetric
		"type": "aes256-gcm96",
	})
	if err == nil {
		//List Keys
		secret, _ = client.Logical().List("/transit/keys")
		printSecret(secret)

		//Read Key Info
		secret, _ = client.Logical().Read("transit/keys/aman")
		printSecret(secret)

		//Edit key
		_, err = client.Logical().Write("/transit/keys/aman/config", map[string]interface{}{
			"deletion_allowed":       true,
			"allow_plaintext_backup": true,
		})
		fmt.Println(err)

		//Rotate Key
		_, err := client.Logical().Write("/transit/keys/aman/rotate", nil)
		secret, _ = client.Logical().Read("transit/keys/aman")
		fmt.Println("Rotated", err)
		printSecret(secret)
		latestVersion := secret.Data["latest_version"]

		//Encrypt Data
		baseData := base64.StdEncoding.EncodeToString([]byte("aman-secret"))
		secret, err := client.Logical().Write("/transit/encrypt/aman", map[string]interface{}{
			"plaintext": baseData,
		})
		fmt.Println("Encrypt", err)
		printSecret(secret)
		cipher := secret.Data["ciphertext"].(string)

		//Export Key
		secret, err = client.Logical().Read("/transit/export/encryption-key/aman/latest")
		fmt.Println("Export Encryption", err)
		keyMap := secret.Data["keys"].(map[string]interface{})
		encryptionKey := keyMap[fmt.Sprintf("%v", latestVersion)].(string)
		decodedEncryptionKey, err := base64.StdEncoding.DecodeString(encryptionKey)
		printSecret(secret)

		secret, err = client.Logical().Read("/transit/export/hmac-key/aman/latest")
		fmt.Println("Export Hmac", err)
		printSecret(secret)

		//Backup Key
		secret, err = client.Logical().Read("/transit/backup/aman")
		fmt.Println("Backup Key", err)
		//backupKey := secret.Data["backup"].(string)
		printSecret(secret)

		//Encrypt-Decrypt Using Key
		text := "aman"
		AAD := []byte("additional authenticated data")
		generatedCipher, noonce, err := dhutil.EncryptAES([]byte(decodedEncryptionKey), []byte(text), AAD)
		plaintext, err := dhutil.DecryptAES([]byte(decodedEncryptionKey), generatedCipher, noonce, AAD)
		fmt.Println("Encrypt/Decrypt GCM:", string(generatedCipher), string(plaintext))

		//Decode using Key
		fmt.Println("Vault Decryption:", cipher)

		//Delete Key
		_, err = client.Logical().Delete("transit/keys/aman")
		fmt.Println("Delete", err)
	}
	return err
}

func printSecret(secret *api.Secret) {
	bytes, _ := json.MarshalIndent(secret, "", "\t")
	fmt.Println(string(bytes))
}

func secretReadWrite(client *api.Client) (err error) {
	var data = map[string]interface{}{
		"Id":       1,
		"Name":     "Aman",
		"password": "Preet",
	}

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

func decrypt(data []byte, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
