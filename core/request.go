package core

import (
	"BlumBot/helper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	username   string
	authToken  string
	httpClient *http.Client
}

func (c *Client) makeRequest(method string, url string, jsonBody interface{}) ([]byte, error) {
	// Convert body to JSON
	var reqBody []byte
	var err error
	if jsonBody != nil {
		reqBody, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	// Create new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	setHeader(req, c.authToken)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 status code
	if resp.StatusCode >= 400 {
		// Read the response body to include in the error message
		bodyBytes, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("error status: %v, and failed to read body: %v", resp.StatusCode, bodyErr)
		}
		return nil, fmt.Errorf("error status: %v, error message: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Login
func (c *Client) getToken(account *Account) map[string]interface{} {
	payload := map[string]string{
		"query": account.QueryData,
	}

	req, err := c.makeRequest("POST", "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to login: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Daily Check in
func (c *Client) dailyCheckIn() string {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/daily-reward?offset=-420", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to check in: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Get Wallet Info
func (c *Client) getWalletInfo() map[string]interface{} {
	req, err := c.makeRequest("GET", "https://wallet-domain.blum.codes/api/v1/wallet/my", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get wallet info: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Get Balance Info
func (c *Client) getBalanceInfo() map[string]interface{} {
	req, err := c.makeRequest("GET", "https://game-domain.blum.codes/api/v1/user/balance", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get balance info: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Get Tribe Info
func (c *Client) getTribeInfo() map[string]interface{} {
	req, err := c.makeRequest("GET", "https://tribe-domain.blum.codes/api/v1/tribe/my", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get tribe info: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Join Tribe
func (c *Client) joinTribe() string {
	res, err := c.makeRequest("POST", "https://tribe-domain.blum.codes/api/v1/tribe/3f4bce0c-9047-4e70-ae12-cb13f92c1196/join", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to join tribe: %v", c.username, err))
	}

	return string(res)
}

// Leave Tribe
func (c *Client) leaveTribe() string {
	var payload map[string]interface{}

	res, err := c.makeRequest("POST", "https://tribe-domain.blum.codes/api/v1/tribe/leave", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to leave tribe: %v", c.username, err))
	}

	return string(res)
}

// Start Farming
func (c *Client) startFarming() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/farming/start", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start farming: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Farming
func (c *Client) claimFarming() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/farming/claim", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim farming: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Task List
func (c *Client) getTask() []map[string]interface{} {
	req, err := c.makeRequest("GET", "https://earn-domain.blum.codes/api/v1/tasks", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get task: %v", c.username, err))
		return nil
	}

	res, err := handleResponseArray(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Start Task
func (c *Client) startTask(taskId string, taskTitle string) map[string]interface{} {
	req, err := c.makeRequest("POST", fmt.Sprintf("https://earn-domain.blum.codes/api/v1/tasks/%s/start", taskId), nil)

	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start task %v: %v", c.username, taskTitle, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Claim Task
func (c *Client) claimTask(taskId string, taskTitle string) map[string]interface{} {
	req, err := c.makeRequest("POST", fmt.Sprintf("https://earn-domain.blum.codes/api/v1/tasks/%s/claim", taskId), nil)

	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to completing task %v: %v", c.username, taskTitle, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Start Game
func (c *Client) startGame() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/game/play", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start game: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Game
func (c *Client) claimGame(gameId string, points int) string {
	payload := map[string]interface{}{
		"gameId": gameId,
		"points": points,
	}

	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/game/claim", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to completing game: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Claim Ref point
func (c *Client) claimRef() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://user-domain.blum.codes/api/v1/friends/claim", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim ref: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// TODO Connect Wallet
func (c *Client) connectWallet(address string, chainId string, publicKey string, payloadId string, signature string, timestamp int) string {
	payload := map[string]interface{}{
		"account": map[string]interface{}{
			"address":   address,
			"chain":     chainId,
			"publicKey": publicKey,
		},
		"tonProof": map[string]interface{}{
			"name": "ton_proof",
			"proof": map[string]interface{}{
				"domain": map[string]interface{}{
					"lengthBytes": 19,
					"value":       "telegram.blum.codes",
				},
				"payload":   payloadId,
				"signature": signature,
				"timestamp": timestamp,
			},
		},
	}

	res, err := c.makeRequest("POST", "https://wallet-domain.blum.codes/api/v1/wallet/connect", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to connect wallet: %v", c.username, err))
		return ""
	}

	return string(res)
}
