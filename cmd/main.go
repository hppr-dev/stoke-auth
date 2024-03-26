package main

import (
	"crypto/ecdsa"
	"log"
	"stoke/internal/cfg"
	"stoke/internal/ctx"
	"stoke/internal/key"
	"stoke/internal/web"
)

func main() {
	log.Println("Starting Stoke Server...")

	config := cfg.Config{
		Server: cfg.Server {
			Address : "",
			Port: 8080,
		},
	}

	tokenIss := &key.AsymetricTokenIssuer[*ecdsa.PrivateKey, *ecdsa.PublicKey]{
		KeyPair: &key.ECDSAKeyPair{ NumBits : 256 },
	}

	context := &ctx.Context{
		Config : config,
		Issuer: tokenIss,
	}

	context.Issuer.Init()
	
	server := web.Server {
		Context: context,
	}

	server.Init()
	if err := server.Run(); err != nil {
		log.Printf("An error stopped the server: %v", err)
	}
	
	log.Println("Stoke Server Terminated.")
}
