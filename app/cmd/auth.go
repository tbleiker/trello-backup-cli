package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var urlAPIKey = "https://trello.com/app-key/"

var urlToken = "https://trello.com/1/authorize?" +
	"expiration=<EXPIRATION>" +
	"&name=<APP_NAME>" +
	"&scope=read" +
	"&response_type=token" +
	"&key=<KEY>"

var cmdAuth = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

func init() {
	cmdAuthToken.Flags().StringVarP(&s.AppName, "app-name", "n", "TrelloBackupCLI", "name of the application")
	cmdAuthToken.Flags().StringVarP(&s.AppExpiration, "app-expiration", "e", "never", "expiration time of the key")
	cmdAuthToken.Flags().StringVarP(&s.APIKey, "api-key", "k", "", "API key")

	cmdAuth.AddCommand(cmdAuthAPIKey)
	cmdAuth.AddCommand(cmdAuthToken)
}

var cmdAuthAPIKey = &cobra.Command{
	Use:   "key",
	Short: "Get API key",
	Run: func(cmd *cobra.Command, args []string) {
		openBrowser(urlAPIKey)
	},
}

var cmdAuthToken = &cobra.Command{
	Use:   "token",
	Short: "Get token",
	Run: func(cmd *cobra.Command, args []string) {
		if s.APIKey == "" {
			exitWithError("API key not specified")
		}

		url := urlToken
		url = strings.Replace(url, "<EXPIRATION>", s.AppExpiration, 1)
		url = strings.Replace(url, "<APP_NAME>", s.AppName, 1)
		url = strings.Replace(url, "<KEY>", s.APIKey, 1)
		openBrowser(url)
	},
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}
