package role

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const ASSIGNABLE_ROLE_NAME_PREFIX = '!'

func isSelfAssignable(roleName string) bool {
	return roleName != "" && roleName[0] == ASSIGNABLE_ROLE_NAME_PREFIX
}

func assignRole(s *discordgo.Session, guildID string, memberID string, role discordgo.Role) error {
	if !isSelfAssignable(role.Name) {
		return errors.New(fmt.Sprintf("Role %s is not self assignable", role.Name))
	}

	return s.GuildMemberRoleAdd(guildID, memberID, role.ID)
}

func removeRole(s *discordgo.Session, guildID string, memberID string, role discordgo.Role) error {
	if !isSelfAssignable(role.Name) {
		return errors.New(fmt.Sprintf("Role %s is not self assignable", role.Name))
	}

	return s.GuildMemberRoleRemove(guildID, memberID, role.ID)
}
