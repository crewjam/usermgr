package cmd

import (
	"io"

	"github.com/crewjam/usermgr"
)

func syncOnce(config *Config, userData *usermgr.UsersData, dryRun bool, stdout io.Writer) error {
	if err := usermgr.SyncUsers(userData, config.LoginGroups, dryRun, stdout); err != nil {
		return err
	}
	if err := usermgr.SyncSudoers(userData, config.LoginGroups, config.SudoGroups, dryRun); err != nil {
		return err
	}
	return nil
}
