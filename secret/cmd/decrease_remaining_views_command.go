package cmd

import (
	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)	

type DecreaseRemainingViewsCommand struct{
	vault secret.Vault
	hash string
}

func NewDecreaseRemainingViewsCommand(vault secret.Vault, hash string) *DecreaseRemainingViewsCommand {
	return &DecreaseRemainingViewsCommand{
		vault: vault,
		hash: hash,
	}
}

func (cmd *DecreaseRemainingViewsCommand) Execute() error {
	if err := cmd.vault.DecreaseRemainingViews(cmd.hash); err != nil {
		return errors.Wrapf(err, "Could not decrease number of views in secret %s", cmd.hash)
	}

	return nil
}
