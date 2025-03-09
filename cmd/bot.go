package main

import (
	"log"
	"mochi-bot/internal/pkg/config"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// ignore error due to .env file not being present
	_ = godotenv.Load()

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	log.Println("Closing session...")

	err = session.Close()
	if err != nil {

		log.Printf("Failed to close session gracefully: %v", err)
	}
}
