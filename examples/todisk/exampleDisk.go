package main

import (
	"log"

	"github.com/ae0000/avatar"
)

func main() {
	if err := avatar.ToDisk("AE", "../../ae.png"); err != nil {
		log.Fatalf("Failed to create AE avatar: %v", err)
	}
	if err := avatar.ToDisk("J", "../../j.png"); err != nil {
		log.Fatalf("Failed to create J avatar: %v", err)
	}
	log.Println("Avatars created successfully!")
}
