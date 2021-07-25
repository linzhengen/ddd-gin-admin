package hash

import (
	//nolint:gosec
	"crypto/md5"
	//nolint:gosec
	"crypto/sha1"
	"fmt"
)

func MD5(b []byte) string {
	//nolint:gosec
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MD5String(s string) string {
	return MD5([]byte(s))
}

func SHA1(b []byte) string {
	//nolint:gosec
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1String(s string) string {
	return SHA1([]byte(s))
}
