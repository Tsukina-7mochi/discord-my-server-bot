package dice

import (
	"fmt"
	"mochi-bot/internal/pkg/botlog"
	"mochi-bot/internal/pkg/discordoptions"

	"github.com/bwmarrin/discordgo"
)

type DiceCommandHandler struct {
	logger *botlog.Logger
}

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

func NewCommandHandler(logger *botlog.Logger) DiceCommandHandler {
	return DiceCommandHandler{
		logger: logger,
	}
}

func (DiceCommandHandler) SubscribingToCommand(name string) bool {
	return name == command.Name
}

func (DiceCommandHandler) NewCommand() discordgo.ApplicationCommand {
	return command
}

func (h DiceCommandHandler) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if rec := recover(); rec != nil {
			h.logger.CommandErrorf(command.Name, "Uncaught error: %v", rec)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "何かしらのエラーが発生しました。",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
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
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
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
		h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
	}

	h.logger.CommandInfof(command.Name, "Rolled %s → %d", dice, result)
}
