package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/dhutil"
)

func transitFun(client *api.Client) (err error) {
	if err == nil {
		//Export Key
		secret, err = client.Logical().Read("/transit/export/encryption-key/aman/latest")
		fmt.Println("Export Encryption", err)
		keyMap := secret.Data["keys"].(map[string]any)

		encryptionKey := keyMap[fmt.Sprintf("%v", latestVersion)].(string)
		decodedEncryptionKey, err := base64.StdEncoding.DecodeString(encryptionKey)

		//Encrypt-Decrypt Using Key
		text := "aman"
		AAD := []byte("additional authenticated data")
		generatedCipher, noonce, err := dhutil.EncryptAES([]byte(decodedEncryptionKey), []byte(text), AAD)
		plaintext, err := dhutil.DecryptAES([]byte(decodedEncryptionKey), generatedCipher, noonce, AAD)
		fmt.Println("Encrypt/Decrypt GCM:", string(generatedCipher), string(plaintext))
	}
	return err
}

func secretReadWrite(client *api.Client) (err error) {
	var data = map[string]any{
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
