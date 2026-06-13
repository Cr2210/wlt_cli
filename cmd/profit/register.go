package profit

import "github.com/spf13/cobra"

// Register adds profit commands to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(profitEventCmd, profitCalculationCmd)
}
