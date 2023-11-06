package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"to-persist/client/global"
	"to-persist/client/util"
)

func init() {
	cobra.OnInitialize(initConfig, util.InitLogger)
	rootCmd.Flags().StringVar(&ConfigFilePath,
		"config",
		"",
		"config file (default is $HOME/.toper/config.yaml)",
	)

}

func initConfig() {

	// Don't forget to read config either from ConfigFilePath or by default!
	if ConfigFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFilePath)
	} else {
		// Use Default config file
		homeDir, _ := os.UserHomeDir()
		ConfigFilePath = filepath.Join(homeDir, ".toper/config.yaml")
		viper.SetConfigFile(ConfigFilePath)

	}

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Panicf("Fatal error config file: %s \n", err.Error())
		// os.Exit(1)
	}

	err := viper.Unmarshal(global.Config)
	if err != nil {
		zap.S().Panicf("Failed to unmarshal global config: %s \n", err.Error())
	}

}

var (
	ConfigFilePath string
	rootCmd        = &cobra.Command{
		Use:   "toper",
		Short: "Toper is a command line tool for keeping track of what you need to persist daily.",
		Args:  cobra.NoArgs,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
