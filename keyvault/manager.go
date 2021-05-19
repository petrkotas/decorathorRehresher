package keyvault

import (
	"time"
)

// Manager manages keyvault
type Manager interface {
	GetSecret(int) *string
}

// manages keyvault
type manager struct {
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) GetSecret(key int) (secret *string) {
	secretStr := RandStringBytes(key)
	secret = &secretStr

	return secret
}

type refreshManager struct {
	interval time.Duration
	manager  Manager
}

func NewRefreshingManager(interval time.Duration, manager Manager) Manager {
	return &refreshManager{
		interval: interval,
		manager:  manager,
	}
}

func (r *refreshManager) GetSecret(key int) *string {
	ticker := time.NewTicker(r.interval)
	done := make(chan bool)

	// create local memory to shield changing memory from GetSecret
	var secret string
	var pSecret *string = &secret

	secret = *r.manager.GetSecret(key)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				secret = *r.manager.GetSecret(key)
			}
		}
	}()

	return pSecret
}
