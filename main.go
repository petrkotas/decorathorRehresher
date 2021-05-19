package main

import (
	"fmt"
	"time"

	"decorator.dev/decorate/keyvault"
)

func main() {
	m := keyvault.NewManager()
	r := keyvault.NewRefreshingManager(1*time.Second, m)

	secret := r.GetSecret(6)

	for i := 0; i < 3; i++ {
		fmt.Println(*secret)
		time.Sleep(3 * time.Second)
	}
}
