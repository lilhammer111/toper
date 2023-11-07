package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"to-persist/client/handler"
)

func init() {
	rootCmd.AddCommand(loginCmd, logoutCmd, registerCmd, userCmd)

	//registerCmd.Flags().StringVar(&global.UserFlags.Mobile, "mobile", "", "Phone number for registration")
	registerCmd.Flags().String("mobile", "", "Phone number for registration")
	err := registerCmd.MarkFlagRequired("mobile")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var (

	// toper register lilhammer111 -m 12312313212
	registerCmd = &cobra.Command{
		Use:   "register <User Name>",
		Short: "Register a new user",
		Long: `The 'register' command allows you to create a new user account for the Toper application.
				Provide the necessary credentials such as username and password 
				to complete the registration process.`,
		Args:   cobra.ExactArgs(1),
		PreRun: handler.RequestToSendSms,
		Run:    handler.Register,
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to your account",
		Long: `The 'login' command lets you access your Toper account.
				You'll need to provide your registered username and password to authenticate and gain access.`,
		Args: cobra.NoArgs,
		Run:  handler.Login,
	}

	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Log out from the current session",
		Long: `The 'logout' command ends your current Toper session, 
				ensuring that no unauthorized actions can be taken on your behalf.
				It's a good practice to log out when you're done using the application, 
				especially on shared machines.`,
		Run: handler.Logout,
	}

	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Show current user information",
		Long: `The 'logout' command ends your current Toper session, 
				ensuring that no unauthorized actions can be taken on your behalf.
				It's a good practice to log out when you're done using the application, 
				especially on shared machines.`,
		Run: handler.Detail,
	}
)
