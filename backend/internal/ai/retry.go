package ai

import "time"

func retry[T any](attempts int, sleep time.Duration, fn func() (T, error)) (T, error) {
	var out T
	var err error
	for i := 0; i < attempts; i++ {
		out, err = fn()
		if err == nil {
			return out, nil
		}
		time.Sleep(sleep)
	}
	return out, err
}
