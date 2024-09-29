package core

import (
	"BlumBot/tools"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func (account *Account) worker(wg *sync.WaitGroup, semaphore *chan struct{}, totalPointsChan *chan int, index int, query string, proxyList []string, walletList []string) {
	defer wg.Done()
	*semaphore <- struct{}{}

	var points int
	var proxy string

	if len(proxyList) > 0 {
		proxy = proxyList[index%len(proxyList)]
	}

	tools.Logger("info", fmt.Sprintf("| %s | Starting Bot...", account.username))

	setDns(&net.Dialer{})

	client := Client{
		account: *account,
		proxy:   proxy,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	if len(client.proxy) > 0 {
		err := client.setProxy()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to set proxy: %v", account.username, err))
		} else {
			tools.Logger("success", fmt.Sprintf("| %s | Proxy Successfully Set...", account.username))
		}
	}

	infoIp, err := client.checkIp()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to check ip: %v", err))
	}

	if infoIp != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Ip: %s | City: %s | Country: %s | Provider: %s", account.username, infoIp["ip"].(string), infoIp["city"].(string), infoIp["country"].(string), infoIp["org"].(string)))
	}

	points = client.autoCompleteTask()

	*totalPointsChan <- points

	<-*semaphore
}

func (c *Client) autoCompleteTask() int {

	var points int

	token, err := c.getToken()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		return 0
	}

	isDailyCheckIn, err := c.dailyCheckIn()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to daily check in: %v", c.account.username, err))
	}

	fmt.Println(isDailyCheckIn)

	walletAddress, err := c.getWalletInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get wallet info: %v", c.account.username, err))
	}

	if walletAddress != "" {
		tools.Logger("success", fmt.Sprintf("| %s | Wallet Address: %s", c.account.username, walletAddress))
	}

	userBalance, err := c.getBalanceInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance info: %v", c.account.username, err))
	}

	if userBalance != nil {
		if balance, exits := userBalance["availableBalance"].(string); exits && len(balance) > 0 {
			if farming, exits := userBalance["farming"].(map[string]interface{}); exits {
				tools.Logger("success", fmt.Sprintf("| %s | Points: | %s | Play Pass: %v | Farming Points: %s", c.account.username, userBalance["availableBalance"].(string), int(userBalance["playPasses"].(float64)), farming["balance"].(string)))
			}
		}
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info balance", c.account.username))
	}

	userTribe, err := c.getTribeInfo()
	if userTribe == nil {
		c.joinTribe()
		time.Sleep(3 * time.Second)
		userTribe, err = c.getTribeInfo()
	} else {
		if tribe, exits := userTribe["id"].(string); exits && tribe != "3f4bce0c-9047-4e70-ae12-cb13f92c1196" {
			c.leaveTribe()
			time.Sleep(3 * time.Second)
			c.joinTribe()
			userTribe, err = c.getTribeInfo()
		}
	}

	if tribeTitle, exits := userTribe["title"].(string); exits && len(tribeTitle) > 0 {
		tools.Logger("success", fmt.Sprintf("| %s | %s | Member: %v | Balance: | %s | Rank: %v", c.account.username, tribeTitle, int(userTribe["countMembers"].(float64)), userTribe["earnBalance"].(string), int(userTribe["rank"].(float64))))
	}

	if farming, exits := userBalance["farming"].(map[string]interface{}); exits {
		if int(farming["endTime"].(float64)) < int(userBalance["timestamp"].(float64)) {
			claimFarming, err := c.claimFarming()
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to claim farming: %v", c.account.username, err))
			}

			if claimFarming != "" {
				tools.Logger("success", fmt.Sprintf("| %s | Claim Farming Successfully | Current Balance: %v | Start Farming After 5s...", c.account.username, claimFarming))

				time.Sleep(5 * time.Second)

				startFarming, err := c.startFarming()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to start farming: %v", c.account.username, err))
				}

				if startFarming != nil {
					if endTime, exits := startFarming["endTime"].(float64); exits {
						claimTime := int(endTime-startFarming["startTime"].(float64)) / 1000
						tools.Logger("success", fmt.Sprintf("| %s | Start Farming Successfully | Claim After: %vs", c.account.username, claimTime))
					}
				} else {
					tools.Logger("error", fmt.Sprintf("| %s | Start Farming Failed", c.account.username))
				}
			} else {
				tools.Logger("error", fmt.Sprintf("| %s | Claim Farming Failed", c.account.username))
			}
		}
	}

	taskList, err := c.getTask()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get task list: %v", c.account.username, err))
	}

	if taskList != nil {
		if len(taskList) > 0 {
			for _, mainTask := range taskList {
				mainTaskMap := mainTask.(map[string]interface{})
				if tasks, exits := mainTaskMap["tasks"].([]interface{}); exits {
					for _, task := range tasks {
						if taskMap, exits := task.(map[string]interface{}); exits && taskMap != nil {
							if taskMap["status"].(string) == "NOT_STARTED" {
								for _, subTask := range taskMap["subTasks"].([]interface{}) {
									subTaskMap := subTask.(map[string]interface{})
									if subTaskMap["status"].(string) != "FINISHED" {
										taskTitle := subTaskMap["title"].(string)

										startTask, err := c.startTask(subTaskMap["id"].(string))
										if err != nil {
											tools.Logger("error", fmt.Sprintf("| %s | Failed to start task %s: %v", c.account.username, taskTitle, err))
										}

										if status, exits := startTask["status"].(string); exits && status == "STARTED" {
											tools.Logger("success", fmt.Sprintf("| %s | Start Task: %v Successfully | Sleep 5s Before Claim Task...", c.account.username, taskTitle))
										}

										time.Sleep(5 * time.Second)

										claimTask, err := c.claimTask(subTaskMap["id"].(string))
										if err != nil {
											tools.Logger("error", fmt.Sprintf("| %s | Failed to claim task %s: %v", c.account.username, taskTitle, err))
										}

										if claimTask != nil {
											if status, exits := claimTask["status"].(string); exits && status == "FINISHED" {
												tools.Logger("success", fmt.Sprintf("| %s | Claim Task: %v Successfully | Reward: | %s | Sleep 5s Before Next Task...", c.account.username, taskTitle, claimTask["reward"].(string)))
											} else {
												tools.Logger("error", fmt.Sprintf("| %s | Claim Task: %v Failed | Sleep 5s Before Next Task...", c.account.username, taskTitle))
											}
										}
									}

									time.Sleep(5 * time.Second)
								}
							}
						}
					}
				}

				if subSections, exits := mainTaskMap["subSections"].([]interface{}); exits {
					for _, sections := range subSections {
						sectionsMap := sections.(map[string]interface{})
						for _, task := range sectionsMap["tasks"].([]interface{}) {
							taskMap := task.(map[string]interface{})
							if taskMap["status"].(string) != "FINISHED" {
								taskTitle := taskMap["title"].(string)
								startTask, err := c.startTask(taskMap["id"].(string))
								if err != nil {
									tools.Logger("error", fmt.Sprintf("| %s | Failed to start task %s: %v", c.account.username, taskTitle, err))
								}

								if status, exits := startTask["status"].(string); exits && status == "STARTED" {
									tools.Logger("success", fmt.Sprintf("| %s | Start Task: %v Successfully | Sleep 5s Before Claim Task...", c.account.username, taskTitle))
								}

								time.Sleep(5 * time.Second)

								claimTask, err := c.claimTask(taskMap["id"].(string))
								if err != nil {
									tools.Logger("error", fmt.Sprintf("| %s | Failed to claim task %s: %v", c.account.username, taskTitle, err))
								}

								if claimTask != nil {
									if status, exits := claimTask["status"].(string); exits && status == "FINISHED" {
										tools.Logger("success", fmt.Sprintf("| %s | Claim Task: %v Successfully | Reward: | %s | Sleep 5s Before Next Task...", c.account.username, taskTitle, claimTask["reward"].(string)))
									} else {
										tools.Logger("error", fmt.Sprintf("| %s | Claim Task: %v Failed | Sleep 5s Before Next Task...", c.account.username, taskTitle))
									}
								}
							}
							time.Sleep(5 * time.Second)
						}
					}
				}
			}
		}
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Failed To Get Task List", c.account.username))
	}

	for i := 0; i < int(userBalance["playPasses"].(float64)); i++ {
		gameId, err := c.startGame()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to start game: %v", c.account.username, err))
		}
		if gameId != "" {
			tools.Logger("success", fmt.Sprintf("| %s | Start Game Successfully | Sleep 30s Before Claim...", c.account.username))

			time.Sleep(30 * time.Second)

			playGame, err := c.claimGame(gameId, tools.RandomNumber(config.Int("GAME_POINTS.MIN"), config.Int("GAME_POINTS.MAX")))
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to claim game: %v", c.account.username, err))
			}

			if playGame != "" {
				if playGame == "OK" {
					tools.Logger("success", fmt.Sprintf("| %s | Claim Game Successfully | Sleep 15s Before Next Game...", c.account.username))
				}
			} else {
				tools.Logger("error", fmt.Sprintf("| %s | Claim Game Failed | Sleep 15s Before Next Game...", c.account.username))
			}
		} else {
			tools.Logger("error", fmt.Sprintf("| %s | Start Game Failed | Sleep 15s Before Next Game...", c.account.username))
			continue
		}

		time.Sleep(15 * time.Second)
	}

	userBalance, err = c.getBalanceInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance info: %v", c.account.username, err))
	}

	if userBalance != nil {
		points, _ = strconv.Atoi(userBalance["availableBalance"].(string))
	}

	return points
}
