package secret

type DecreaseRemainingViewsCommand struct{}

func (cmd *DecreaseRemainingViewsCommand) Execute(vault Vault, hash string) error {
	if err := vault.DecreaseRemainingViews(hash); err != nil {
		// todo! errwrapf
		return  err
	}

	return nil
}
