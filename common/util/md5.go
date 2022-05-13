package util

import (
	"crypto/md5"
	"encoding/hex"
)

type Md5Info struct {
	Hash     string
	FileList []string
	Count    int
}

func (self *Md5Info) Add(path string) {
	self.FileList = append(self.FileList, path)
	self.Count++
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
