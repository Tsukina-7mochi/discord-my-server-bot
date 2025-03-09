package role

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const ASSIGNABLE_ROLE_NAME_PREFIX = '#'

func isSelfAssignable(s *discordgo.Session, guildID string, roleID string) (bool, error) {
	role, err := s.State.Role(guildID, roleID)
	if err != nil {
		return false, err
	}

	return role.Name[0] == ASSIGNABLE_ROLE_NAME_PREFIX, nil
}

func assignRole(s *discordgo.Session, guildID string, memberID string, roleID string) error {
	isSelfAssignable, err := isSelfAssignable(s, guildID, roleID)
	if err != nil {
		return err
	}

	if !isSelfAssignable {
		return errors.New(fmt.Sprintf("Role %s is not self assignable", roleID))
	}

	return s.GuildMemberRoleAdd(guildID, memberID, roleID)
}

func removeRole(s *discordgo.Session, guildID string, memberID string, roleID string) error {
	isSelfAssignable, err := isSelfAssignable(s, guildID, roleID)
	if err != nil {
		return err
	}

	if !isSelfAssignable {
		return errors.New(fmt.Sprintf("Role %s is not self assignable", roleID))
	}

	return s.GuildMemberRoleRemove(guildID, memberID, roleID)
}
