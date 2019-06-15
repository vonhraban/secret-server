package secret

type DecreaseRemainingViewsCommand struct{
	vault Vault
	hash string
}

func NewDecreaseRemainingViewsCommand(vault Vault, hash string) *DecreaseRemainingViewsCommand {
	return &DecreaseRemainingViewsCommand{
		vault: vault,
		hash: hash,
	}
}

func (cmd *DecreaseRemainingViewsCommand) Execute() error {
	if err := cmd.vault.DecreaseRemainingViews(cmd.hash); err != nil {
		// todo! errwrapf
		return  err
	}

	return nil
}
