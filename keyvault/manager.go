package keyvault

import (
	"time"
)

type certificateGetter func(int) []*string

// Manager manages keyvault
type Manager interface {
	GetSecret(int) []*string
}

// manages keyvault
type manager struct {
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) GetSecret(key int) (secret []*string) {
	secret = make([]*string, 3)
	for i := 0; i < 3; i++ {
		secretStr := RandStringBytes(key)
		secret[i] = &secretStr
	}

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

func (r *refreshManager) GetSecret(key int) []*string {
	ticker := time.NewTicker(r.interval)
	done := make(chan bool)

	tmpSecret := r.manager.GetSecret(key)

	// create local memory to shield changing memory from GetSecret
	var secret []string = make([]string, len(tmpSecret))
	var pSecret []*string = make([]*string, len(tmpSecret))

	for i, v := range tmpSecret {
		secret[i] = *v
	}

	for i := 0; i < len(secret); i++ {
		pSecret[i] = &secret[i]
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for i, v := range r.manager.GetSecret(key) {
					secret[i] = *v
				}
			}
		}
	}()

	return pSecret
}

func RefreshingCerticates(interval time.Duration, done chan bool, getter certificateGetter) certificateGetter {
	var secret []string
	var pSecret []*string
	certs := func(key int) []*string {
		ticker := time.NewTicker(interval)
		tmpSecret := getter(key)

		secret = make([]string, len(tmpSecret), 10)
		pSecret = make([]*string, len(tmpSecret), 10)

		// initial memory copy, to shield changing memory from GetSecret
		for i, v := range tmpSecret {
			secret[i] = *v
		}
		for i := 0; i < len(secret); i++ {
			pSecret[i] = &secret[i]
		}

		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					for i, v := range getter(key) {
						secret[i] = *v
					}
				}
			}
		}()

		return pSecret
	}

	return certs
}
