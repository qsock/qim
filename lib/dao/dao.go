package dao

import (
	"sync"
)

var (
	env string
	mu  sync.Mutex
	kmu sync.Mutex
)

func SetEnv(_env string) {
	env = _env
}
