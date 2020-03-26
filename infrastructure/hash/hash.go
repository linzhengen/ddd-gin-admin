package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

// MD5Hash creates md5 hash.
func MD5Hash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5HashString md5 hash string
func MD5HashString(s string) string {
	return MD5Hash([]byte(s))
}

// SHA1Hash create SHA1 hash.
func SHA1Hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1HashString SHA1 hash string.
func SHA1HashString(s string) string {
	return SHA1Hash([]byte(s))
}
