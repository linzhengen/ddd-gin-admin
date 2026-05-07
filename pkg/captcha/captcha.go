// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package captcha implements generation and verification of image CAPTCHAs.
package captcha

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/linzhengen/ddd-gin-admin/pkg/captcha/store"
)

const (
	// DefaultLen Default number of digits in captcha solution.
	DefaultLen = 6
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
)

var (
	// ErrNotFound id not found
	ErrNotFound = errors.New("captcha: id not found")

	globalStore store.Store
	once        sync.Once
)

func getStore() store.Store {
	once.Do(func() {
		if globalStore == nil {
			globalStore = store.NewMemoryStore(time.Second, Expiration)
		}
	})
	return globalStore
}

// SetCustomStore sets custom storage for captchas, replacing the default
// memory store. This function must be called before generating any captchas.
func SetCustomStore(s store.Store) {
	globalStore = s
}

// New creates a new captcha with the standard length, saves it in the internal
// storage and returns its id.
func New() string {
	return NewLen(DefaultLen)
}

// NewLen is just like New, but accepts length of a captcha solution as the
// argument.
func NewLen(length int) (id string) {
	id = randomID()
	getStore().Set(id, RandomDigits(length))
	return
}

// Reload generates and remembers new digits for the given captcha id.  This
// function returns false if there is no captcha with the given id.
func Reload(id string) bool {
	old := getStore().Get(id, false)
	if old == nil {
		return false
	}
	getStore().Set(id, RandomDigits(len(old)))
	return true
}

// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func WriteImage(w io.Writer, id string, width, height int) error {
	d := getStore().Get(id, false)
	if d == nil {
		return ErrNotFound
	}
	_, err := NewImage(id, d, width, height).WriteTo(w)
	return err
}

// Verify returns true if the given digits are the ones that were used to
// create the given captcha id.
func Verify(id string, digits []byte) bool {
	if digits == nil || len(digits) == 0 {
		return false
	}
	reald := getStore().Get(id, true)
	if reald == nil {
		return false
	}
	return bytes.Equal(digits, reald)
}

// VerifyString is like Verify, but accepts a string of digits.  It removes
// spaces and commas from the string.
func VerifyString(id string, digits string) bool {
	if digits == "" {
		return false
	}
	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
			// ignore
		default:
			return false
		}
	}
	return Verify(id, ns)
}
