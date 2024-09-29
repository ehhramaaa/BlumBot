package core

import (
	"BlumBot/tools"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func (account *Account) parsingQueryData() {
	value, err := url.ParseQuery(account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse query data: %s", err))
	}

	if len(value.Get("query_id")) > 0 {
		account.queryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.authDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.hash = value.Get("hash")
	}

	userParam := value.Get("user")

	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse user data: %s", err))
	}

	userId, ok := userData["id"].(float64)
	if !ok {
		tools.Logger("error", "Failed to convert ID to float64")
	}

	account.userId = int(userId)

	username, ok := userData["username"].(string)
	if !ok {
		tools.Logger("error", "Failed to get username from query")
		return
	}

	account.username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get first name from query")
	}

	account.firstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get last name from query")
	}
	account.lastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		tools.Logger("error", "Failed to get language code from query")
	}
	account.languageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		tools.Logger("error", "Failed to get allows write to pm from query")
	}

	account.allowWriteToPm = allowWriteToPm
}

func LaunchBot() {
	queryPath := "configs/query.txt"
	proxyPath := "configs/proxy.txt"
	maxThread := config.Int("MAX_THREAD")
	isUseProxy := config.Bool("USE_PROXY")

	queryData, err := tools.ReadFileTxt(queryPath)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Query Data Not Found: %s", err))
		return
	}

	tools.Logger("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	var wg sync.WaitGroup
	var semaphore chan struct{}
	var proxyList, walletList []string

	if isUseProxy {
		proxyList, err = tools.ReadFileTxt(proxyPath)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("Proxy Data Not Found: %s", err))
		}

		tools.Logger("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))
	}

	totalPointsChan := make(chan int, len(queryData))

	if maxThread > len(queryData) {
		semaphore = make(chan struct{}, len(queryData))
	} else {
		semaphore = make(chan struct{}, maxThread)
	}

	for {
		for index, query := range queryData {
			wg.Add(1)
			account := &Account{
				queryData: query,
			}

			account.parsingQueryData()

			go account.worker(&wg, &semaphore, &totalPointsChan, index, query, proxyList, walletList)
		}
		wg.Wait()
		close(totalPointsChan)

		var totalPoints int

		for points := range totalPointsChan {
			totalPoints += points
		}

		tools.Logger("success", fmt.Sprintf("Total Points All Account: %v", totalPoints))

		randomSleep := tools.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))

		tools.Logger("info", fmt.Sprintf("Launch Bot Finished | Sleep %vs Before Next Lap...", randomSleep))

		time.Sleep(time.Duration(randomSleep) * time.Second)
	}
}
