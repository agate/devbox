// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package boxcli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/devbox/internal/devbox"
	"go.jetpack.io/devbox/internal/devbox/devopt"
	envsecIntegration "go.jetpack.io/devbox/internal/integrations/envsec"
)

type envsecInitCmdFlags struct {
	config configFlags
	force  bool
}

func envsecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envsec",
		Short: "envsec commands",
	}
	cmd.AddCommand(envsecInitCmd())
	cmd.Hidden = true
	return cmd
}

func envsecInitCmd() *cobra.Command {
	flags := envsecInitCmdFlags{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize envsec integration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return envsecInitFunc(cmd, flags)
		},
	}

	flags.config.register(cmd)
	cmd.Flags().BoolVarP(
		&flags.force,
		"force",
		"f",
		false,
		"Force initialization even if already initialized",
	)

	return cmd
}

func envsecInitFunc(cmd *cobra.Command, flags envsecInitCmdFlags) error {
	box, err := devbox.Open(&devopt.Opts{
		Dir:    flags.config.path,
		Stderr: cmd.ErrOrStderr(),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	envsec := envsecIntegration.DefaultEnvsec(cmd.ErrOrStderr(), box.ProjectDir())
	if err := envsec.NewProject(cmd.Context(), flags.force); err != nil {
		return errors.WithStack(err)
	}
	box.Config().SetStringField("EnvFrom", "envsec")
	return box.Config().SaveTo(box.ProjectDir())
}
