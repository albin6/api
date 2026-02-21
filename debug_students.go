package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BaseURL        = "https://api.brototype.com/tool/api"
	Mobile         = "+91-9539547616"
	Password       = "ponnuAamy@11"
	TransactionID  = "1b75f425-01d9-4c7e-8159-4b011e092090"
	HardcodedToken = "0cAFcWeA54HDfc_jI8MCGdzrRFlDHaofZhj03RB46Ox6Gd8ZYJg_xe9ZKxk2UX5JXPP8ejC9nwf2YhTpGDgzs_O_WcOfRFdgnXfKungZRaaNeEfbLvkU4cIGa6hKKcVBwHiFKT1YwxtGRwrhrBzlOywAl3HGr7A-d11mNReXA-fcUiwKn948iC8wqtrQ72lX3n8ih8mMUEfyB9G13gI6G58oFqBf8_kyuQ_6GOyASWfHU1aBpgDIlAD5LTbNkIB6mZ0Vqyehn1X6qfcsiV_CzaMx9yxz9iXfY4kxov09W-34mVK6LdvAz4mHPCsrQxk4xxtUeXbMSl8QTeYfWC63C7FCoXp1jbmuK8A9lUzhJNWw6rqzaIwUrCewEo7Jbi1HkwpzDT_DKQy_Qnfvpkm5fxgHhi6Iz1_98b5-VJNpDBt8KNgRYAJgI8rPe-qcJMBvj8Krz1WTYt5gSr6HJKLDSNMpm8J4Edlb9WfrBWX8jcOtK_bY41GTHb0VA82eDXXPiQ_E9DozdmjmGFusmGE0KgptFFpzWL0DKR0kwkda2iKe_GQhMlk4stBfmKZnPvOhqhYLnLGJtfN9RSCZqBymvDrKq22GbrHLSDWy5NUj1gU6lCQFZGF10m8zhlbt_ZacxwnY1J-eHUihzp9DjnW0qQGR97abssgJGj7snTzjcmlrBkBb-Z6exFHk57UUkPwVrapdxe6d76G7Ru-9Y4vNnPGmJvqWmGDp0kKcrzN3w4iWLfnAiZR-8xaw6Xr2g926cjiGuQ_TPXlSPs7Uq7H0wQf8CVSnmNpmA74LUE8p5yFYU42wTv_j7KtVSegPAwaOVhbFCcpWKz2vqIfxJmYEi5baPyvZFCqPrFmWptHd0x7YaiR9UJlwRr2b-Pz_M96sofUn5xmUn18l9NQSDHQz5rrrLh1NnqLg2y5LAJ9hqV_715QJAaIz5Hn2QFtTSU54Srz5ofp7HaTXEG32RIo3fJemIG3F1aK-jNjiKxdv773lBG4b1SEv0zAhtQmd3Nth_vIW2OhBswn0JUcOs2rbd0hmZLIcpfaU1Cy0azwjNkcci2FFiXh3lAofexZU3SjJUhZZAHMpBZC_xg9iV9sTH3TqIQg2wQBaLjkNKV_DjLYlUfhAwij7_7qWVP6dCNDo_zO2m9_Zlyy1agb7InrJzETF_ozaaAtRVKGeOR91iVOSTDl-i10_uPuioGTmNQFtTaft69s_vGVEfUo69IM2I9aPz5gYMAUFmx1ZApobwEUhSLaMbe2j4X0zKAIOYTZdoHyIK3Mr5UtTAPVX5r0bCSCGo_520cXlfz9oJJgabyDQxXD1glngCpCHDQWLhWMz0rqE5DQ-mMCbmVOXB5n99j4F4AP0CxJB6TXzmskzU91XTsHIC6UkNQ-uQah4U5jYPju3fVyEsF7i9DblALnE2QJdUFwo2MvLwADoySgLG12lTyO_bKBtgwDpV9wKiTJ-Et5Jv7Bklao8Bk8SMlrxupj_7qvb8NCEpTBquwNrTIPRFaeerCe1Z4iZ1kP156Gx_Co2QIBjNl8FRYazIQEC5CUlmYTd3ZVrgdIpoO7zcOyQIGp38T5_xWPM6kcFPVfc7gKc5UOKxLeHZBt_xBJzXKgQUO6UUjSfj4bFnGFBWeAqvsJAz1OrBz1JAh4tRpvfH2FClLBD_mfvc9zrpSYdOGi9yeSWwV9jMl94lO_5Z-tjxBHY9NfaKHaM7uZzzk5Ws4VnnMaTvirLZGYE8Df2IAVYnB8qyvLq4-yOkeeYPuwr3Un96ClahA3DRg884UnOP8wA8Ks5g2rhpHMxBqQX3et-hfcZt8r547ObXSB0EYDXypJgyqp8QMGU0gKfG0sc5fgCihaYl6zSMn6HJm7hxzrF3ZCUv7XNHnxBfSsT_4INu843jVTKdDU6FWNhUzbR78_dV816slcdWZYIzolW6IUSC7Kh63Uzr2d4sPFQULfLkbmqU3QfjqPmx08DUGLuTdrTADu3UdJBwrUslCIm7eVis4ihTiXEjoQHSkKl3Ody73p1wDpZV4WCT4ZbB1W4Mh--rgkJweA-G_iGJv5qUjW-2AMQqI1SDDMcXuSkmY2BlLvi12O0lwK3crEJU-hS8Mu7r8mOJHsjYJVe8OZqbsRm6L-f9ndezjhQ"
)

type AuthRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func main() {
	// Step 1: Subscribe/Authenticate
	fmt.Println("Authenticating...")
	token := HardcodedToken

	// Check if hardcoded token works, otherwise try to authenticate
	// For now let's try to authenticate explicitly to be sure
	authBody := AuthRequest{
		Mobile:   Mobile,
		Password: Password,
		Token:    HardcodedToken,
	}
	jsonData, _ := json.Marshal(authBody)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", BaseURL+"/auth/verifySignIn", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-transaction-id", TransactionID)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Auth failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var authResp AuthResponse
		if err := json.Unmarshal(bodyBytes, &authResp); err == nil && authResp.Data.Token != "" {
			token = authResp.Data.Token
			fmt.Println("Got new token.")
		}
	} else {
		fmt.Println("Auth failed or using hardcoded token.")
	}

	// Step 2: Get Students with Status=2
	fmt.Println("\nFetching Students (status=2)...")
	req, _ = http.NewRequest("GET", BaseURL+"/student", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-transaction-id", TransactionID)

	q := req.URL.Query()
	q.Add("statusCodes", "2")
	q.Add("offset", "0")
	q.Add("pageSize", "10")
	q.Add("transactionId", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req.URL.RawQuery = q.Encode()

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Get Students failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ = io.ReadAll(resp.Body)
	fmt.Printf("Get Students Status: %d\n", resp.StatusCode)

	// Write to file
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
		os.WriteFile("debug_students.json", prettyJSON.Bytes(), 0644)
		fmt.Println("Wrote response to debug_students.json")
	} else {
		fmt.Printf("Response Body (Raw): %s\n", string(bodyBytes))
	}
}
