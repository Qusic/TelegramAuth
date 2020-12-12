package main

import "time"

func runTask() {
	timeoutTick := time.Tick(config.authTimeout)
	for {
		select {
		case <-timeoutTick:
			checkExpiration(time.Now())
		}
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
