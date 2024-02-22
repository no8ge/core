package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "core",
		Short: "A container controler service for atop in k8s.",
		Long:  `A container controler service for atop in k8s.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("mod", "m", "release", "mod of core service")
	rootCmd.PersistentFlags().StringP("port", "p", "8080", "port of core service")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.core.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(wd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".core")
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed using config file: ", err)
	}
}
