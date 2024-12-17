package git

import (
	"context"
	"os/exec"
)

func ShallowClone(ctx context.Context, repository string, output string, directories ...string) error {
	cloneCmd := exec.CommandContext(ctx, "git", "clone", "--filter=tree:0", "--depth=1", "--no-checkout", "--sparse", repository, output)
	if err := cloneCmd.Run(); err != nil {
		return err
	}

	sparseInitCmd := exec.CommandContext(ctx, "git", "sparse-checkout", "init", "--sparse-index", "--cone")
	sparseInitCmd.Dir = output
	if err := sparseInitCmd.Run(); err != nil {
		return err
	}

	sparseInitOptions := append([]string{
		"sparse-checkout",
		"add",
	}, directories...)
	spareInitCheckoutCmd := exec.CommandContext(ctx, "git", sparseInitOptions...)
	spareInitCheckoutCmd.Dir = output
	if err := spareInitCheckoutCmd.Run(); err != nil {
		return err
	}

	checkoutCmd := exec.CommandContext(ctx, "git", "checkout")
	checkoutCmd.Dir = output
	if err := checkoutCmd.Run(); err != nil {
		return err
	}

	return nil
}
