package dice

import (
	"fmt"
	"log"
	"mochi-bot/internal/pkg/discordoptions"

	"github.com/bwmarrin/discordgo"
)

type DiceCommandHandler struct{}

var command = discordgo.ApplicationCommand{
	Name:        "dice",
	Description: "ダイスを振る",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "dice",
			Description: "振るダイス (例: 2d6)",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func NewCommandHandler() DiceCommandHandler {
	return DiceCommandHandler{}
}

func (DiceCommandHandler) SubscribingToCommand(name string) bool {
	return name == command.Name
}

func (DiceCommandHandler) NewCommand() discordgo.ApplicationCommand {
	return command
}

func (DiceCommandHandler) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Failed to handle command %s: %v", i.ID, rec)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something weng wrong :(",
					Flags:   1 << 6, // ephemeral
				},
			})
			if err != nil {
				log.Printf("[Command %s] Failed to respond interaction: %v", i.ID, err)
				return
			}
		}
	}()

	data := i.ApplicationCommandData()
	if data.Name != command.Name {
		return
	}
	options := discordoptions.ParseOptions(data.Options)
	dice := options["dice"].StringValue()
	diceRoll := parseToDiceRoll(dice)
	if diceRoll == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ダイスのフォーマットが正しくありません！(例: 2d6)",
				Flags:   1 << 6, // ephemeral
			},
		})
		if err != nil {
			log.Printf("[Command %s] Failed to respond interaction: %v", i.ID, err)
			return
		}
	}

	result := diceRoll.Roll()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s -> %d", dice, result),
		},
	})
	if err != nil {
		log.Printf("[Command %s] Failed to respond interaction: %v", i.ID, err)
	}

	log.Printf("[Dice] Rolled %s → %d", dice, result)
}
