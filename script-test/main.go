package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/muhadif/e-wallet/core/entity"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type DepositPayloadTest struct {
	Amount     int32
	DelayAfter int32 //in second
}

// TestScenario Please run service before run this test
func main() {
	baseUrl := "http://0.0.0.0:8000"

	walletTest1 := getRandomWallet(10)
	walletTest2 := getRandomWallet(10)
	walletTest3 := getRandomWallet(10)
	walletTest4 := getRandomWallet(10)
	tests := []struct {
		name           string
		want           string
		payloadDeposit []*DepositPayloadTest
		walletID       string
		wantErr        bool
	}{
		{
			name: "given two deposits of 6,000 amount each, both within 2 minutes should return amount 12000 and IsAboveThreshold=true",
			want: "{\"status\":\"success\",\"data\":{\"WalletID\":\"%s\",\"Amount\":12000,\"IsAboveThreshold\":true}}\n",
			payloadDeposit: []*DepositPayloadTest{
				{
					Amount:     6000,
					DelayAfter: 10,
				},
				{
					Amount:     6000,
					DelayAfter: 0,
				},
			},
			walletID: walletTest1,
		},
		{
			name: "given one single deposit of 6,000, then after 2-minutes later another single deposit of\n6,000 should return amount 12000 and IsAboveThreshold=false",
			want: "{\"status\":\"success\",\"data\":{\"WalletID\":\"%s\",\"Amount\":12000,\"IsAboveThreshold\":false}}\n",
			payloadDeposit: []*DepositPayloadTest{
				{
					Amount:     6000,
					DelayAfter: 121,
				},
				{
					Amount:     6000,
					DelayAfter: 0,
				},
			},
			walletID: walletTest2,
		},
		{
			name: "given five deposits of 2,000 amount each all within 2 minutes, then after 5 seconds later\nanother single deposit of 6,000 should return amount 16,000 and IsAboveThreshold=false",
			want: "{\"status\":\"success\",\"data\":{\"WalletID\":\"%s\",\"Amount\":16000,\"IsAboveThreshold\":false}}\n",
			payloadDeposit: []*DepositPayloadTest{
				{
					Amount:     2000,
					DelayAfter: 24,
				},
				{
					Amount:     2000,
					DelayAfter: 24,
				},
				{
					Amount:     2000,
					DelayAfter: 24,
				},
				{
					Amount:     2000,
					DelayAfter: 24,
				},
				{
					Amount:     2000,
					DelayAfter: 24,
				},
				{
					Amount:     6000,
					DelayAfter: 12,
				},
			},
			walletID: walletTest3,
		},
		{
			name: "six deposits of 2,000 amount each all within 2 minutes should return amount 12000 and IsAboveThreshold=true",
			want: "{\"status\":\"success\",\"data\":{\"WalletID\":\"%s\",\"Amount\":12000,\"IsAboveThreshold\":true}}\n",
			payloadDeposit: []*DepositPayloadTest{
				{
					Amount:     2000,
					DelayAfter: 20,
				},
				{
					Amount:     2000,
					DelayAfter: 20,
				},
				{
					Amount:     2000,
					DelayAfter: 20,
				},
				{
					Amount:     2000,
					DelayAfter: 20,
				},
				{
					Amount:     2000,
					DelayAfter: 20,
				},
				{
					Amount:     2000,
					DelayAfter: 20,
				},
			},
			walletID: walletTest4,
		},
	}
	for _, tt := range tests {
		log.Println(fmt.Sprintf("Run test for %s", tt.name))
		if err := depositTestScenarioWithPayload(baseUrl, tt.payloadDeposit, tt.walletID); err != nil {
			log.Fatal(fmt.Errorf("TestScenario() error = %v, wantErr %v", err, tt.wantErr))
		}

		// wait data sync
		time.Sleep(1 * time.Second)
		got, err := checkDepositAmount(baseUrl, tt.walletID)
		if err != nil {
			log.Fatal(fmt.Errorf("TestScenario() error = %v, wantErr %v", err, tt.wantErr))
		}

		want := fmt.Sprintf(tt.want, tt.walletID)
		if got != want {
			log.Fatal(fmt.Errorf("TestScenario() got = %v, want %v", got, want))
		}

		log.Print(fmt.Sprintf("Got test result %s", got))
		log.Print(fmt.Sprintf("Want test result%s", want))
		log.Println(fmt.Sprintf("Finish test for %s", tt.name))
		log.Println()
	}
}

func checkDepositAmount(baseUrl string, walletID string) (string, error) {
	walletUrl := fmt.Sprintf("/api/wallet/%s", walletID)
	request, err := http.NewRequest("GET", baseUrl+walletUrl, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error : %s", response.Body)
	}
	body, _ := ioutil.ReadAll(response.Body)

	return string(body), nil
}

func depositTestScenarioWithPayload(baseUrl string, payloads []*DepositPayloadTest, walletID string) error {
	payloadLength := len(payloads)
	for idx, payloadTest := range payloads {
		payload := &entity.DepositByWalletID{
			WalletID:      walletID,
			DepositAmount: int(payloadTest.Amount),
		}
		payload.WalletID = walletID
		payloadByte, _ := json.Marshal(payload)

		request, err := http.NewRequest("POST", baseUrl+"/api/deposit", bytes.NewBuffer(payloadByte))
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			return fmt.Errorf("error : %s", response.Body)
		}
		//body, _ := ioutil.ReadAll(response.Body)
		log.Println(fmt.Sprintf("success deposit with payload : %s  --> delay : %d second for the next", string(payloadByte), payloadTest.DelayAfter))

		if idx < payloadLength-1 {
			delayAfter := time.Duration(payloadTest.DelayAfter)
			time.Sleep(delayAfter * time.Second)
		}
	}

	return nil
}

func getRandomWallet(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:n]
}
