package esyncdefine

import "time"

var defaultRetryDelay = time.Second

func SetDefaultRetryDelay(delay time.Duration) {
	defaultRetryDelay = delay
}

func GetDefaultRetryDelay() time.Duration {
	return defaultRetryDelay
}
