package utils

import (
	"crypto/md5"
	"encoding/hex"
)

type md5Lib struct{}

var Md5Lib = md5Lib{}

func (md5Lib) MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func (lib md5Lib) Md5AndSalt(password string, salt string) string {
	return lib.MD5(password + salt)
}
