package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var urlList = "https://api.trello.com/1/members/me/boards" +
	"?key=<KEY>" +
	"&token=<TOKEN>"

var urlGet = "https://api.trello.com/1/boards/<ID>" +
	"?key=<KEY>" +
	"&token=<TOKEN>" +
	"&actions=all" +
	"&actions_limit=1000" +
	"&cards=all" +
	"&lists=all" +
	"&members=all" +
	"&member_fields=all" +
	"&checklists=all" +
	"&fields=all"

type boards struct {
	Name string
	ID   string
}

func init() {
	cmdList.Flags().StringVarP(&s.APIKey, "api-key", "k", "", "API key")
	cmdList.Flags().StringVarP(&s.Token, "token", "t", "", "token")

	cmdGet.Flags().StringVarP(&s.APIKey, "api-key", "k", "", "API key")
	cmdGet.Flags().StringVarP(&s.Token, "token", "t", "", "token")
	cmdGet.Flags().StringVarP(&s.BoardName, "board-name", "n", "", "Board name(s)")
	cmdGet.Flags().StringVarP(&s.BoardID, "board-id", "i", "", "Board ID(s)")
	cmdGet.Flags().StringVarP(&s.BackupPath, "backup-path", "p", "", "Path to store backup")
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List available boards",
	Run: func(cmd *cobra.Command, args []string) {
		url := getURL(urlList)
		boards := getBoards(url)

		for _, b := range boards {
			fmt.Printf("%s: %s\n", b.ID, b.Name)

		}
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get board(s)",
	Run: func(cmd *cobra.Command, args []string) {
		urlList := getURL(urlList)
		boards := getBoards(urlList)

		if s.BoardID == "" {
			if s.BoardName == "" {
				exitWithError("Board ID or name not specified")
			}

			for _, b := range boards {
				if b.Name == s.BoardName {
					cmd.Flags().Set("board-id", fmt.Sprintf("%v", b.ID))
					break
				}
			}

			if s.BoardID == "" {
				exitWithError(fmt.Sprintf("Board '%s' does not exist", s.BoardName))
			}
		} else {
			for _, b := range boards {
				if b.ID == s.BoardID {
					cmd.Flags().Set("board-name", fmt.Sprintf("%v", b.Name))
					break
				}
			}

			if s.BoardName == "" {
				exitWithError(fmt.Sprintf("Board with ID '%s' does not exist", s.BoardID))
			}
		}

		urlGet := getURL(urlGet)
		res, err := http.Get(urlGet)
		if err != nil {
			exitWithError(err.Error())
		}

		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()

		var fileName string
		var filePath string

		time := time.Now()
		fileName = fmt.Sprintf("trello_%s_%s.json", s.BoardName, time.Format("20060102-150405"))

		if s.BackupPath == "" {
			filePath = s.AppPath
		} else {
			filePath = s.BackupPath
		}
		file := filepath.Join(filePath, fileName)

		f, err := os.Create(file)
		if err != nil {
			exitWithError(err.Error())
		}
		defer f.Close()

		_, err = f.WriteString(string(data))
		if err != nil {
			exitWithError(err.Error())
		}
		log.Info(fmt.Sprintf("backup written to %s", file))
	},
}

func getURL(url string) string {
	if s.APIKey == "" {
		exitWithError("API key not specified")
	}
	if s.Token == "" {
		exitWithError("Token not specified")
	}

	url = strings.Replace(url, "<ID>", s.BoardID, 1)
	url = strings.Replace(url, "<KEY>", s.APIKey, 1)
	url = strings.Replace(url, "<TOKEN>", s.Token, 1)

	return url
}

func getBoards(url string) []boards {
	var err error
	var b []boards

	res, err := http.Get(url)
	if err != nil {
		exitWithError(err.Error())
	}

	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		exitWithError(err.Error())
	}

	return b
}
