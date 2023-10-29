package command

import (
	"github.com/spf13/cobra"
	"to-persist/client/handler"
)

func init() {
	rootCmd.AddCommand(versionCmd, pingCmd)
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Toper",
		Long:  `All software has versions. This is Toper's`,
		Run:   handler.ViewVersion,
	}
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Test network connectivity",
		Long: `The 'ping' command is used to test the network connectivity 
				between the host and a specified address.
				It sends a series of packets to the target address and awaits a response. 
				This can be useful for diagnosing network issues 
				and checking the availability of a remote server.`,
		Run: handler.Ping,
	}
)
