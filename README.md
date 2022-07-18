# E-Wallet Test

# Test Case Result
```
2022/07/19 01:16:23 Run test for given two deposits of 6,000 amount each, both within 2 minutes should return amount 12000 and IsAboveThreshold=true
2022/07/19 01:16:23 success deposit with payload : {"walletID":"7da5aabf6b","depositAmount":6000}  --> delay : 10 second for the next
2022/07/19 01:16:33 success deposit with payload : {"walletID":"7da5aabf6b","depositAmount":6000}  --> delay : 0 second for the next
2022/07/19 01:16:43 Got test result {"status":"success","data":{"WalletID":"7da5aabf6b","Amount":12000,"IsAboveThreshold":true}}
2022/07/19 01:16:43 Want test result{"status":"success","data":{"WalletID":"7da5aabf6b","Amount":12000,"IsAboveThreshold":true}}
2022/07/19 01:16:43 Finish test for given two deposits of 6,000 amount each, both within 2 minutes should return amount 12000 and IsAboveThreshold=true
2022/07/19 01:16:43 
2022/07/19 01:16:43 Run test for given one single deposit of 6,000, then after 2-minutes later another single deposit of
6,000 should return amount 12000 and IsAboveThreshold=false
2022/07/19 01:16:43 success deposit with payload : {"walletID":"be9d10fde2","depositAmount":6000}  --> delay : 121 second for the next
2022/07/19 01:18:44 success deposit with payload : {"walletID":"be9d10fde2","depositAmount":6000}  --> delay : 0 second for the next
2022/07/19 01:18:54 Got test result {"status":"success","data":{"WalletID":"be9d10fde2","Amount":12000,"IsAboveThreshold":false}}
2022/07/19 01:18:54 Want test result{"status":"success","data":{"WalletID":"be9d10fde2","Amount":12000,"IsAboveThreshold":false}}
2022/07/19 01:18:54 Finish test for given one single deposit of 6,000, then after 2-minutes later another single deposit of
6,000 should return amount 12000 and IsAboveThreshold=false
2022/07/19 01:18:54 
2022/07/19 01:18:54 Run test for given five deposits of 2,000 amount each all within 2 minutes, then after 5 seconds later
another single deposit of 6,000 should return amount 16,000 and IsAboveThreshold=false
2022/07/19 01:18:54 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":2000}  --> delay : 24 second for the next
2022/07/19 01:19:19 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":2000}  --> delay : 24 second for the next
2022/07/19 01:19:43 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":2000}  --> delay : 24 second for the next
2022/07/19 01:20:07 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":2000}  --> delay : 24 second for the next
2022/07/19 01:20:31 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":2000}  --> delay : 24 second for the next
2022/07/19 01:20:55 success deposit with payload : {"walletID":"3c6a1dbbe2","depositAmount":6000}  --> delay : 12 second for the next
2022/07/19 01:21:05 Got test result {"status":"success","data":{"WalletID":"3c6a1dbbe2","Amount":16000,"IsAboveThreshold":false}}
2022/07/19 01:21:05 Want test result{"status":"success","data":{"WalletID":"3c6a1dbbe2","Amount":16000,"IsAboveThreshold":false}}
2022/07/19 01:21:05 Finish test for given five deposits of 2,000 amount each all within 2 minutes, then after 5 seconds later
another single deposit of 6,000 should return amount 16,000 and IsAboveThreshold=false
2022/07/19 01:21:05 
2022/07/19 01:21:05 Run test for six deposits of 2,000 amount each all within 2 minutes should return amount 12000 and IsAboveThreshold=true
2022/07/19 01:21:05 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:21:25 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:21:45 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:22:06 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:22:26 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:22:46 success deposit with payload : {"walletID":"251fbb132b","depositAmount":2000}  --> delay : 20 second for the next
2022/07/19 01:22:56 Got test result {"status":"success","data":{"WalletID":"251fbb132b","Amount":12000,"IsAboveThreshold":true}}
2022/07/19 01:22:56 Want test result{"status":"success","data":{"WalletID":"251fbb132b","Amount":12000,"IsAboveThreshold":true}}
2022/07/19 01:22:56 Finish test for six deposits of 2,000 amount each all within 2 minutes should return amount 12000 and IsAboveThreshold=true

```

# Setup Environment
1. Run kafka broker and any other tools docker-compose.yml in project using following command
```
docker compose up
```
2. Create kafka topic for development : deposits, balance-table, above-threshold-table\
   You can edit topic in ```config/config.go```
```
   TopicDeposit   		goka.Stream = "deposits"
	GroupBalance   		goka.Group  = "balance"
	GroupAboveThreshold goka.Group = "above-threshold"
```
3. Make sure broker host is up and configured right in ```config/config.go```
```
var Brokers = []string{"localhost:29092"}
```
3. After all requirement ok, move to the installations

# Installation
1. Install library that needed for development
```
go mod tidy
go mod download
```
or using
```
go get
```
2. Generate proto with following command
```
cd proto
protoc --go_out=.  *.proto
```

3. Run core service using following command
```
go run main.go
```
3. Make sure no error and the terminal shown below
```                                    
2022/07/18 00:23:12 All component is running OK, enjoy!
2022/07/18 00:23:15 [Processor balance2] setup generation 84, claims=map[string][]int32{"depo1":[]int32{0}}
2022/07/18 00:23:15 [Processor above-threshold-3] setup generation 21, claims=map[string][]int32{"depo1":[]int32{0}}
```
4. Service is running now!

## Script Test (Automatic Test based on test case)
1. Make sure service is running in port 8000 (you can modify port in test file)
2. Run script with following command
```
go run ./script-test/main.go
```

# Manual Test (Hit endpoint)
1. Endpoint ```/api/deposit```, run with command bellow
```
curl --location --request GET '0.0.0.0:8000/api/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "walletID": "6bdd384452",
    "depositAmount": 10000
}'
```
2. Endpoint ```/api/wallet/{{walletID}}```, run with command bellow
```
curl --location --request GET '0.0.0.0:8000/api/wallet/6bdd384452'
```