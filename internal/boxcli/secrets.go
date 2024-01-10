// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package boxcli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/devbox/internal/devbox"
	"go.jetpack.io/devbox/internal/devbox/devopt"
	"go.jetpack.io/envsec/pkg/envsec"
)

type secretsFlags struct {
	config configFlags
}

type secretsInitCmdFlags struct {
	force bool
}

type secretsListFlags struct {
	show   bool
	format string
}

func secretsCmd() *cobra.Command {
	flags := &secretsFlags{}
	cmd := &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"envsec"},
		Short:   "Interact with devbox secrets.",
	}
	cmd.AddCommand(secretsInitCmd(flags))
	cmd.AddCommand(secretsListCmd(flags))
	cmd.AddCommand(secretsSetCmd(flags))
	cmd.Hidden = true

	flags.config.register(cmd)

	return cmd
}

func secretsInitCmd(secretsFlags *secretsFlags) *cobra.Command {
	flags := secretsInitCmdFlags{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize secrets management with jetpack cloud",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretsInitFunc(cmd, flags, secretsFlags)
		},
	}

	cmd.Flags().BoolVarP(
		&flags.force,
		"force",
		"f",
		false,
		"Force initialization even if already initialized",
	)

	return cmd
}

func secretsSetCmd(flags *secretsFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "set <NAME1>=<value1> [<NAME2>=<value2>]...",
		Short: "Securely store one or more environment variables",
		Long:  "Securely store one or more environment variables. To test contents of a file as a secret use set=@<file>",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return envsec.ValidateSetArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			box, err := devbox.Open(&devopt.Opts{
				Dir:    flags.config.path,
				Stderr: cmd.ErrOrStderr(),
			})
			if err != nil {
				return errors.WithStack(err)
			}

			secrets, err := box.Secrets(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			envID, err := secrets.EnvID()
			if err != nil {
				return errors.WithStack(err)
			}

			return secrets.SetFromArgs(ctx, envID, args)
		},
	}
}

func secretsListCmd(commonFlags *secretsFlags) *cobra.Command {
	flags := secretsListFlags{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all secrets",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			box, err := devbox.Open(&devopt.Opts{
				Dir:    commonFlags.config.path,
				Stderr: cmd.ErrOrStderr(),
			})
			if err != nil {
				return errors.WithStack(err)
			}

			secrets, err := box.Secrets(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			envID, err := secrets.EnvID()
			if err != nil {
				return errors.WithStack(err)
			}

			vars, err := secrets.List(ctx, envID)
			if err != nil {
				return err
			}

			return envsec.PrintEnvVars(
				vars, cmd.OutOrStdout(), flags.show, flags.format)
		},
	}

	cmd.Flags().BoolVarP(
		&flags.show,
		"show",
		"s",
		false,
		"Display the value of each environment variable (secrets included)",
	)
	cmd.Flags().StringVarP(
		&flags.format,
		"format",
		"f",
		"table",
		"Display the key values in key=value format. Must be one of: table | dotenv | json",
	)
	return cmd
}

func secretsInitFunc(
	cmd *cobra.Command,
	flags secretsInitCmdFlags,
	secretsFlags *secretsFlags,
) error {
	ctx := cmd.Context()
	box, err := devbox.Open(&devopt.Opts{
		Dir:    secretsFlags.config.path,
		Stderr: cmd.ErrOrStderr(),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	secrets, err := box.Secrets(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := secrets.NewProject(ctx, flags.force); err != nil {
		return errors.WithStack(err)
	}
	box.Config().SetStringField("EnvFrom", "envsec")
	return box.Config().SaveTo(box.ProjectDir())
}
