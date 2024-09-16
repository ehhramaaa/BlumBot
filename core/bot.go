package core

import (
	"BlumBot/helper"
	"fmt"
	"net/http"
	"time"
)

func launchBot(account *Account, gamePoints int, walletAddress string) {
	client := &Client{
		username:   account.Username,
		httpClient: &http.Client{},
	}

	tokens := client.getToken(account)

	if token, exits := tokens["token"].(map[string]interface{}); exits {
		client.authToken = fmt.Sprintf("Bearer %s", token["access"].(string))
	} else {
		helper.PrettyLog("error", "Failed To Get Token")
		return
	}

	isDailyCheckIn := client.dailyCheckIn()
	if isDailyCheckIn == "OK" {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Daily Check in Successfully...", client.username))
	}

	userWallet := client.getWalletInfo()
	if address, exits := userWallet["address"].(string); exits && len(address) > 0 {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Wallet Address: %s", client.username, address))
	}

	userBalance := client.getBalanceInfo()
	if balance, exits := userBalance["availableBalance"].(string); exits && len(balance) > 0 {
		if farming, exits := userBalance["farming"].(map[string]interface{}); exits {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Points: | %s | Play Pass: %v | Farming Points: %s", client.username, userBalance["availableBalance"].(string), int(userBalance["playPasses"].(float64)), farming["balance"].(string)))
		}
	}

	userTribe := client.getTribeInfo()
	if userTribe == nil {
		client.joinTribe()
		time.Sleep(3 * time.Second)
		userTribe = client.getTribeInfo()
	}

	if tribe, exits := userTribe["id"].(string); exits && tribe != "3f4bce0c-9047-4e70-ae12-cb13f92c1196" {
		leaveTribe := client.leaveTribe()
		time.Sleep(3 * time.Second)
		if leaveTribe == "OK" {
			time.Sleep(3 * time.Second)
			client.joinTribe()
			time.Sleep(3 * time.Second)
			userTribe = client.getTribeInfo()
		}
	}

	if tribeTitle, exits := userTribe["title"].(string); exits && len(tribeTitle) > 0 {
		helper.PrettyLog("success", fmt.Sprintf("| %s | | %s | Member: %v | Balance: | %s | Rank: %v", client.username, tribeTitle, int(userTribe["countMembers"].(float64)), userTribe["earnBalance"].(string), int(userTribe["rank"].(float64))))
	}

	if farming, exits := userBalance["farming"].(map[string]interface{}); exits {
		if int(farming["endTime"].(float64)) < int(userBalance["timestamp"].(float64)) {
			claimFarming := client.claimFarming()
			if currentBalance, exits := claimFarming["availableBalance"].(string); exits && len(currentBalance) > 0 {
				helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Farming Successfully | Current Balance: %v | Start Farming After 5s...", client.username, currentBalance))

				time.Sleep(5 * time.Second)

				startFarming := client.startFarming()
				if endTime, exits := startFarming["endTime"].(float64); exits {
					claimTime := int(endTime-startFarming["startTime"].(float64)) / 1000
					helper.PrettyLog("success", fmt.Sprintf("| %s | Start Farming Successfully | Claim After: %vs", client.username, claimTime))
				}
			}
		}
	}

	taskList := client.getTask()

	if len(taskList) > 0 {
		for _, mainTask := range taskList {
			if tasks, exits := mainTask["tasks"].([]interface{}); exits {
				for _, task := range tasks {
					if taskMap, exits := task.(map[string]interface{}); exits && taskMap != nil {
						if taskMap["status"].(string) == "NOT_STARTED" {
							for _, subTask := range taskMap["subTasks"].([]interface{}) {
								subTaskMap := subTask.(map[string]interface{})
								if subTaskMap["status"].(string) != "FINISHED" {
									startTask := client.startTask(subTaskMap["id"].(string), subTaskMap["title"].(string))
									if status, exits := startTask["status"].(string); exits && status == "STARTED" {
										helper.PrettyLog("success", fmt.Sprintf("| %s | Start Task: %v Successfully | Sleep 5s Before Claim Task...", client.username, startTask["title"].(string)))
									}

									time.Sleep(5 * time.Second)

									claimTask := client.claimTask(subTaskMap["id"].(string), subTaskMap["title"].(string))
									if claimTask != nil {
										if status, exits := claimTask["status"].(string); exits && status == "FINISHED" {
											helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Task: %v Successfully | Reward: | %s | Sleep 15s Before Next Task...", client.username, claimTask["reward"].(string), claimTask["title"].(string)))
										} else {
											helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Task: %v Failed | Sleep 15s Before Next Task...", client.username, claimTask["title"].(string)))
										}
									}
								}

								time.Sleep(15 * time.Second)
							}
						}
					}
				}
			}

			if subSections, exits := mainTask["subSections"].([]interface{}); exits {
				for _, sections := range subSections {
					sectionsMap := sections.(map[string]interface{})
					for _, task := range sectionsMap["tasks"].([]interface{}) {
						taskMap := task.(map[string]interface{})
						if taskMap["status"].(string) != "FINISHED" {
							startTask := client.startTask(taskMap["id"].(string), taskMap["title"].(string))
							if status, exits := startTask["status"].(string); exits && status == "STARTED" {
								helper.PrettyLog("success", fmt.Sprintf("| %s | Start Task: %v Successfully | Sleep 5s Before Claim Task...", client.username, startTask["title"].(string)))
							}

							time.Sleep(5 * time.Second)

							claimTask := client.claimTask(taskMap["id"].(string), taskMap["title"].(string))
							if claimTask != nil {
								if status, exits := claimTask["status"].(string); exits && status == "FINISHED" {
									helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Task: %v Successfully | Reward: | %s | Sleep 15s Before Next Task...", client.username, claimTask["reward"].(string), claimTask["title"].(string)))
								} else {
									helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Task: %v Failed | Sleep 15s Before Next Task...", client.username, claimTask["title"].(string)))
								}
							}
						}
						time.Sleep(15 * time.Second)
					}
				}
			}
		}
	}

	for i := 0; i < int(userBalance["playPasses"].(float64)); i++ {
		startGame := client.startGame()
		if gameId, exits := startGame["gameId"].(string); exits && len(gameId) > 0 {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Start Game Successfully | Sleep 30s Before Claim...", client.username))

			time.Sleep(30 * time.Second)

			playGame := client.claimGame(gameId, gamePoints)
			if playGame == "OK" {
				helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Game Successfully | Sleep 15s Before Next Game...", client.username))
			}
		}

		time.Sleep(15 * time.Second)
	}
}
