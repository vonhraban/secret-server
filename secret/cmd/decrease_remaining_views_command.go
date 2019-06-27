package cmd

import (
	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)

type decreaseRemainingViewsCommand struct {
	vault secret.Vault
	hash  string
}

func NewDecreaseRemainingViewsCommand(vault secret.Vault, hash string) *decreaseRemainingViewsCommand {
	return &decreaseRemainingViewsCommand{
		vault: vault,
		hash:  hash,
	}
}

func (cmd *decreaseRemainingViewsCommand) Execute() error {
	if err := cmd.vault.DecreaseRemainingViews(cmd.hash); err != nil {
		return errors.Wrapf(err, "Could not decrease number of views in secret %s", cmd.hash)
	}

	return nil
}
