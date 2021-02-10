package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type settings struct {
	ConfigFile    string
	Verbose       bool
	APIKey        string
	Token         string
	AppName       string
	AppExpiration string
	BoardName     string
	BoardID       string
	BackupPath    string
	AppPath       string
}

var s settings

const (
	configFilename = "trello-backup.yaml"
	envPrefix      = "TRELLO"
)

var rootCmd = &cobra.Command{
	Use:   "trello-backup-cli",
	Short: "Command line tool to backup trello boards.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ex, err := os.Executable()
		if err != nil {
			exitWithError(err.Error())
		}
		s.AppPath = filepath.Dir(ex)

		return initializeConfig(cmd)
	},
}

// Execute command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initializeApp)

	rootCmd.PersistentFlags().StringVarP(&s.ConfigFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&s.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(cmdAuth)
	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdGet)
}

func initializeApp() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
	if s.Verbose == false {
		log.SetLevel(log.WarnLevel)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	var err error

	if s.ConfigFile != "" {
		viper.SetConfigFile(s.ConfigFile)
	} else {
		viper.AddConfigPath(s.AppPath)
		viper.SetConfigFile(configFilename)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err == nil {
		log.Info("Using config file: ", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("trello")
	viper.AutomaticEnv()

	bindFlags(cmd)

	return nil
}

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func exitWithError(msg string) {
	log.Error(msg)
	os.Exit(1)
}
