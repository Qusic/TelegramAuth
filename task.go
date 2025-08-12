package main

import "time"

func runTask() {
	for now := range time.Tick(config.authTimeout) {
		checkExpiration(now)
	}
}

func checkExpiration(now time.Time) {
	state.authMutex.Lock()
	defer state.authMutex.Unlock()
	for token, usedTime := range state.authCache {
		recentlyUsed := now.Sub(usedTime) < config.authTimeout
		if !recentlyUsed {
			delete(state.authCache, token)
		}
	}
}
