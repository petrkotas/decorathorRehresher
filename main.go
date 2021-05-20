package main

import (
	"fmt"
	"time"

	"decorator.dev/decorate/keyvault"
)

func main() {
	m := keyvault.NewManager()
	done := make(chan bool)

	getSecret := keyvault.RefreshingCerticates(1*time.Second, done, m.GetSecret)
	secret := getSecret(6)

	for i := 0; i < 3; i++ {
		for i := 0; i < 3; i++ {
			fmt.Print(*secret[i], ", ")
		}
		fmt.Println()
		time.Sleep(3 * time.Second)
	}
	done <- true
}
