package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckExpiration(t *testing.T) {
	now,
		_, _,
		validToken1, validToken2,
		_, _, _, _ := setupTestToken()
	for _, data := range []struct {
		offset int
		valid1 bool
		valid2 bool
	}{
		{offset: 5, valid1: true, valid2: true},
		{offset: 6, valid1: false, valid2: true},
		{offset: 7, valid1: false, valid2: true},
		{offset: 8, valid1: false, valid2: false},
		{offset: 9, valid1: false, valid2: false},
	} {
		config.authDuration = 5 * time.Second
		config.authTimeout = 5 * time.Second
		state.authCache = map[string]time.Time{}
		useToken(validToken1, now.Add(1*time.Second))
		useToken(validToken2, now.Add(3*time.Second))
		checkExpiration(now.Add(time.Duration(data.offset) * time.Second))
		if data.valid1 {
			assert.Contains(t, state.authCache, validToken1)
		} else {
			assert.NotContains(t, state.authCache, validToken1)
		}
		if data.valid2 {
			assert.Contains(t, state.authCache, validToken2)
		} else {
			assert.NotContains(t, state.authCache, validToken2)
		}
	}
}
