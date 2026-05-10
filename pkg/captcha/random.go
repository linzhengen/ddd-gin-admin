// Copyright 2011-2014 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"sync"
)

// idLen is a length of captcha id string.
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// rngKey is a secret key used to deterministically derive seeds for
// PRNGs used in image. Generated once during initialization.
var (
	rngKey     [32]byte
	rngKeyOnce sync.Once
)

func setRNGKey() {
	rngKeyOnce.Do(func() {
		if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
			panic("captcha: error reading random source: " + err.Error())
		}
	})
}

// Purposes for seed derivation.
const (
	imageSeedPurpose = 0x01
)

// deriveSeed returns a 16-byte PRNG seed from rngKey, purpose, id and digits.
func deriveSeed(purpose byte, id string, digits []byte) (out [16]byte) {
	setRNGKey()
	var buf [sha256.Size]byte
	h := hmac.New(sha256.New, rngKey[:])
	_, _ = h.Write([]byte{purpose})
	_, _ = io.WriteString(h, id)
	_, _ = h.Write([]byte{0})
	_, _ = h.Write(digits)
	sum := h.Sum(buf[:0])
	copy(out[:], sum)
	return
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9.
func RandomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte((256%int(mod))&0xFF)
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}
}

// randomID returns a new random id string.
func randomID() string {
	b := randomBytesMod(idLen, byte(len(idChars)&0xFF))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}
