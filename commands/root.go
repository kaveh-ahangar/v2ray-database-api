package commands

import (
	"fmt"
	"os"
	"v2ray-database-api/config"

	"github.com/spf13/cobra"
)

var persistentOpts = config.CliOnlyOptions{}
var ConfigPath string
var devMode bool
var rootCmd = &cobra.Command{
	Short: "v2ray api",
	Args:  ValidateRootArgs,
	Use:   "v2ray-database-api",
}

func init() {
	// get config path
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "config file (default is $HOME/.config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&devMode, "dev", "d", false, "enable development mode")
}

// ValidateRootArgs validates the arguments passed to the root command.
func ValidateRootArgs(cmd *cobra.Command, args []string) error {
	// the user must specify at least one argument OR wait for input on stdin IF it is a pipe
	if len(args) == 0 && !IsPipedInput() {
		// return an error with no message for the user, which will implicitly show the help text (but no specific error)
		return fmt.Errorf("")
	}

	return cobra.MaximumNArgs(1)(cmd, args)
}

// IsPipedInput returns true if there is no input device, which means the user **may** be providing input via a pipe.
func IsPipedInput() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeNamedPipe != 0
}
