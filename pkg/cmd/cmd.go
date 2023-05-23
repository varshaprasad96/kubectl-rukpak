package cmd

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-rukpak",
		Short: "Manage installation of operator content on cluster",
	}
	cmd.AddCommand()
	return cmd
}
