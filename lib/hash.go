/**
 * 项目用到的hash算法
 */
package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func MD5(source string) string {
	h := md5.New()
	h.Write([]byte(source))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func MD5FromBytes(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func SHA1(source string) string {
	h := sha1.New()
	h.Write([]byte(source))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func SHA1FromBytes(data []byte) string {
	h := sha1.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
