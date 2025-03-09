package main

import (
	"log"
	"mochi-bot/internal/app/dice"
	"mochi-bot/internal/app/role"
	"mochi-bot/internal/pkg/botlog"
	"mochi-bot/internal/pkg/config"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type CommandHandler interface {
	SubscribingToCommand(string) bool
	NewCommand() discordgo.ApplicationCommand
	HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func createCommands(session *discordgo.Session, appID string, guildID string, handlers []CommandHandler) error {
	commands := make([]*discordgo.ApplicationCommand, 0, len(handlers))
	for _, handler := range handlers {
		command := handler.NewCommand()
		commands = append(commands, &command)
	}
	_, err := session.ApplicationCommandBulkOverwrite(appID, guildID, commands)
	return err
}

func deleteAllCommands(session *discordgo.Session, appID string, guildID string) error {
	commands, err := session.ApplicationCommands(appID, guildID)
	if err != nil {
		return err
	}

	for _, command := range commands {
		err = session.ApplicationCommandDelete(appID, guildID, command.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// ignore error due to .env file not being present
	_ = godotenv.Load()

	config, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: Failed to load config: %v", err)
	}

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("FATAL: Failed to create session: %v", err)
	}

	logger := botlog.NewLogger(log.Default(), session, config.LogChannelID)
	commandHandlers := []CommandHandler{
		dice.NewCommandHandler(logger),
		role.NewCommandHandler(logger),
	}

	log.Printf("INFO: Deleting commands")
	err = deleteAllCommands(session, config.AppID, config.GuildID)
	if err != nil {
		log.Fatalf("FATAL: Failed to delete commands: %v", err)
	}

	log.Printf("INFO: Creating commands")
	err = createCommands(session, config.AppID, config.GuildID, commandHandlers)
	if err != nil {
		log.Fatalf("FATAL: Failed to create commands: %v", err)
	}

	// add a handler to log when the bot is ready
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Infof("Logged in as %s", r.User.String())
	})

	// add a handler to respond to command interactions
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name
		for _, handler := range commandHandlers {
			if handler.SubscribingToCommand(name) {
				handler.HandleCommand(s, i)
				return
			}
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("FATAL: Failed to open session: %v", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	log.Println("INFO: Closing session...")

	err = session.Close()
	if err != nil {
		log.Printf("ERROR: Failed to close session gracefully: %v", err)
	}
}
