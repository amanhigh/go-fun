package util

// nolint:gosec
import (
	"crypto/md5"
	"encoding/hex"
)

type Md5Info struct {
	Hash     string
	FileList []string
	Count    int
}

func (m *Md5Info) Add(path string) {
	m.FileList = append(m.FileList, path)
	m.Count++
}

func GetMD5Hash(text string) string {
	//nolint:gosec
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
