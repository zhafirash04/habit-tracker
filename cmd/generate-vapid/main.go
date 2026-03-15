package main

import (
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalf("Failed to generate VAPID keys: %v", err)
	}

	fmt.Println("=== VAPID Keys Generated ===")
	fmt.Println()
	fmt.Println("Add these to your .env file:")
	fmt.Println()
	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", publicKey)
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", privateKey)
	fmt.Printf("VAPID_SUBJECT=mailto:admin@habitflow.app\n")
	fmt.Println()
	fmt.Println("The PUBLIC key is shared with the browser (frontend).")
	fmt.Println("The PRIVATE key must be kept SECRET (backend only).")
}
