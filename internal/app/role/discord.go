package role

import (
	"mochi-bot/internal/pkg/botlog"
	"mochi-bot/internal/pkg/discordoptions"

	"github.com/bwmarrin/discordgo"
)

type RoleCommandHandler struct {
	logger *botlog.Logger
}

var command = discordgo.ApplicationCommand{
	Name:        "role",
	Description: "ロールの管理",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "assign",
			Description: "ロールを追加",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "role",
					Description: "ロール",
					Type:        discordgo.ApplicationCommandOptionRole,
					Required:    true,
				},
			},
		},
		{
			Name:        "remove",
			Description: "ロールを削除",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "role",
					Description: "ロール",
					Type:        discordgo.ApplicationCommandOptionRole,
					Required:    true,
				},
			},
		},
	},
}

func NewCommandHandler(logger *botlog.Logger) RoleCommandHandler {
	return RoleCommandHandler{
		logger: logger,
	}
}

func (RoleCommandHandler) SubscribingToCommand(name string) bool {
	return name == command.Name
}

func (RoleCommandHandler) NewCommand() discordgo.ApplicationCommand {
	return command
}

func (h RoleCommandHandler) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	operation := data.Options[0].Name

	subOptions := discordoptions.ParseOptions(data.Options[0].Options)
	user := i.Member.User
	role := subOptions["role"].RoleValue(s, i.GuildID)

	if operation == "assign" {
		err := assignRole(s, i.GuildID, user.ID, *role)
		if err != nil {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "このロールは追加できません",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
			}
			return
		}

		h.logger.Infof("Role assigned <@%s>(%s) to <@&%s>(%s)", user.ID, user.Username, role.ID, role.Name)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ロールを追加しました！",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
		}
	} else if operation == "remove" {
		err := removeRole(s, i.GuildID, user.ID, *role)
		if err != nil {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "このロールは削除できません",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
			}
			return
		}

		h.logger.Infof("Role removed <@%s>(%s) to <@&%s>(%s)", user.ID, user.Username, role.ID, role.Name)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ロールを削除しました！",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			h.logger.CommandErrorf(command.Name, "Failed to respond to interaction: %v", err)
		}
	} else {
		panic("Invalid operation")
	}
}
