package core

import "fmt"

// Login
func (c *Client) getToken() (string, error) {
	payload := map[string]string{
		"query": c.account.queryData,
	}

	res, err := c.makeRequest("POST", "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP", payload)
	if err != nil {
		return "", err
	}

	if token, exits := res["token"].(map[string]interface{}); exits {
		return token["access"].(string), nil
	} else {
		return "", fmt.Errorf("Token field not found!")
	}
}

// Daily Check in
func (c *Client) dailyCheckIn() (string, error) {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/daily-reward?offset=-420", nil)
	if err != nil {
		return "", err
	}

	return res["response"].(string), nil
}

// Get Wallet Info
func (c *Client) getWalletInfo() (string, error) {
	res, err := c.makeRequest("GET", "https://wallet-domain.blum.codes/api/v1/wallet/my", nil)
	if err != nil {
		return "", err
	}

	if wallet, exits := res["address"].(string); exits {
		return wallet, nil
	} else {
		return "", fmt.Errorf("Wallet field not found!")
	}
}

// Get Balance Info
func (c *Client) getBalanceInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", "https://game-domain.blum.codes/api/v1/user/balance", nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get Tribe Info
func (c *Client) getTribeInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", "https://tribe-domain.blum.codes/api/v1/tribe/my", nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Join Tribe
func (c *Client) joinTribe() {
	c.makeRequest("POST", "https://tribe-domain.blum.codes/api/v1/tribe/3f4bce0c-9047-4e70-ae12-cb13f92c1196/join", nil)
}

// Leave Tribe
func (c *Client) leaveTribe() {
	var payload map[string]interface{}

	c.makeRequest("POST", "https://tribe-domain.blum.codes/api/v1/tribe/leave", payload)
}

// Start Farming
func (c *Client) startFarming() (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/farming/start", nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Claim Farming
func (c *Client) claimFarming() (string, error) {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/farming/claim", nil)
	if err != nil {
		return "", err
	}

	if balance, exits := res["availableBalance"].(string); exits {
		return balance, nil
	} else {
		return "", nil
	}
}

// Get Task List
func (c *Client) getTask() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", "https://earn-domain.blum.codes/api/v1/tasks", nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Start Task
func (c *Client) startTask(taskId string) (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", fmt.Sprintf("https://earn-domain.blum.codes/api/v1/tasks/%s/start", taskId), nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Claim Task
func (c *Client) claimTask(taskId string) (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", fmt.Sprintf("https://earn-domain.blum.codes/api/v1/tasks/%s/claim", taskId), nil)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Start Game
func (c *Client) startGame() (string, error) {
	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/game/play", nil)
	if err != nil {
		return "", err
	}

	if gameId, exits := res["gameId"].(string); exits {
		return gameId, nil
	} else {
		return "", nil
	}
}

// Claim Game
func (c *Client) claimGame(gameId string, points int) (string, error) {
	payload := map[string]interface{}{
		"gameId": gameId,
		"points": points,
	}

	res, err := c.makeRequest("POST", "https://game-domain.blum.codes/api/v1/game/claim", payload)
	if err != nil {
		return "", err
	}

	return res["response"].(string), nil
}

// Claim Ref point
func (c *Client) claimRef() (string, error) {
	res, err := c.makeRequest("POST", "https://user-domain.blum.codes/api/v1/friends/claim", nil)
	if err != nil {
		return "", err
	}

	if amount, exits := res["claimBalance"].(string); exits {
		return amount, nil
	} else {
		return "", nil
	}
}
