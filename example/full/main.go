package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/guarilha/go-gnosispay"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	client, err := gnosispay.New(nil,
		gnosispay.SetBaseURL("https://api.gnosispay.com"),
		gnosispay.SetSIWEParams("https://your-app.com"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pk, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		fmt.Printf("failed to parse private key: %v", err)
		return
	}

	address := common.HexToAddress(os.Getenv("ADDRESS"))

	// Auth
	_, err = client.Auth.AuthenticateWithPrivateKey(ctx, address, pk)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// User
	user, err := client.User.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Printf("User email: %s\n", user.Email)
	if user.KycStatus != nil {
		fmt.Printf("KYC Status: %s\n", *user.KycStatus)
	}

	// Cards
	cards, err := client.Cards.List(ctx)
	if err != nil {
		log.Fatalf("Failed to get cards: %v", err)
	}
	fmt.Printf("cards: %v\n", cards)

	// KYC
	integration, err := client.KYC.GetIntegration(ctx)
	if err != nil {
		log.Fatalf("Failed to get integration: %v", err)
	}
	fmt.Printf("integration: %v\n", integration)

	// IBAN
	available, err := client.IBAN.CheckAvailability(ctx)
	if err != nil {
		log.Fatalf("Failed to get available: %v", err)
	}
	fmt.Printf("available: %v\n", available)

	// Account
	balances, err := client.Account.GetBalances(ctx)
	if err != nil {
		log.Fatalf("Failed to get balances: %v", err)
	}
	fmt.Printf("balances: %v\n", balances)
}
