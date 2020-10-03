package main

import "github.com/s-zer0/wallet/pkg/wallet"

func main() {
	svc := &wallet.Service{}
	wallet.RegisterAccount(svc, "+992000000001")
}