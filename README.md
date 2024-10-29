# Overview
This Wallet Service allows users to manage their digital wallet by providing functionalities to deposit money, withdraw money, send money to other users, check wallet balance, and view transaction history. The service is implemented as a RESTful API, making it easy to integrate with various front-end applications or other services.

# Features
1. Deposit Money: Users can add funds to their wallet.
2. Withdraw Money: Users can withdraw funds from their wallet.
3. Send Money: Users can transfer money to another user’s wallet.
4. Check Balance: Users can view their current wallet balance.
5. Transaction History: Users can view a log of all their transactions.

# Code structure
I use gin framework, and use `github.com/gin-gonic/gin` to handle http request. The code structure is base on MVC structure, it's very easy to understand and write code.

```
├─cmd              # main file
├─config           # config file
├─controller       # controller
├─data             # data
├─model            # db model
├─router           # gin router
├─service          # service
│  └─dao           # dao layer
└─util             # utils
    ├─db           # db/redis init
    └─errcode      # define error code
```
# Requirements
- golang 1.23.0（or later versions）
- postgresql
- redis

# How to Run
**1. Install postgresql and redis**
1) down load postgresql and redis, run them, create database and tables.

2) install postgresql client(pgadmin4), and connect to postgresql.

3) excute sql script `simplewallet/schema.sql` to create database "wallet" and tables "wallets"、"transactions".

4) start redis service example:
```bash
> sudo service redis-server start
```

**2. Fill config file**

config the `simplewallet/conf.yaml` with your own postgresql and redis config. 

example:
```yaml
env: local
gin_host: :8080
db:
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: 123456
  dbname: wallet
redis:
  uri: 127.0.0.1:6379
  password: 123456
  db: 0
```

**3. Run the service**
```bash
> cd simplewallet
> go mod tidy
> go run cmd/main.go
```

# API Documentation
after program running, use postman or other tools to test the api.Example:

1) POST  http://127.0.0.1:8080/deposit

input param:
```json
{
    "order_id":"111",
    "user_id": 101,
    "amount": 1000.00
}
```

output:
```json
{
    "code": 0,
    "message": "Deposit successful",
    "log_id": "6720d14400032aa0"
}
```

2) POST  http://127.0.0.1:8080/withdraw

input param:
```json
{
    "order_id":"113",
    "user_id": 101,
    "amount": 500.00
}

```

output:
```json
{
    "code": 0,
    "message": "Withdrawal successful",
    "log_id": "6720d3160006df74"
}
```

3) POST  http://127.0.0.1:8080/transfer

input param:
```json
{
    "order_id": "1001",
    "from_user_id": 101,
    "to_user_id": 102,
    "amount": 1000.00
}
```

output:
```json
{
    "code": 0,
    "message": "Transfer successful",
    "log_id": "6720d3d6000a399c"
}
```

4) GET  http://127.0.0.1:8080/balance?user_id=101

output:
```json
{
    "code": 0,
    "message": "Success",
    "data": {
        "balance": 500
    },
    "log_id": "6720d41400080850"
}
```

5) GET  http://127.0.0.1:8080/transactions?user_id=101&page=1&size=10
output:
```json
{
    "code": 0,
    "message": "Success",
    "data": {
        "items": [
            {
                "order_id": "111",
                "user_id": 101,
                "tx_type": 1,
                "amount": 1000,
                "related_user_id": 0,
                "created_at": "2024-10-29 20:12:52"
            },
            {
                "order_id": "112",
                "user_id": 101,
                "tx_type": 1,
                "amount": 1000,
                "related_user_id": 0,
                "created_at": "2024-10-29 20:17:53"
            },
            {
                "order_id": "113",
                "user_id": 101,
                "tx_type": 2,
                "amount": 500,
                "related_user_id": 0,
                "created_at": "2024-10-29 20:20:38"
            },
            {
                "order_id": "1001",
                "user_id": 101,
                "tx_type": 4,
                "amount": 1000,
                "related_user_id": 102,
                "created_at": "2024-10-29 20:23:50"
            }
        ]
    },
    "log_id": "6720d45500030664"
}
```

# Testing
**1. golangci-lint** 
config see `.golangci.yml`.
```
$ golangci-lint run
<testsuites></testsuites>
```

**2. unit test conver**
```
$ go test ./... -race -cover
?       simplewallet/model      [no test files]
?       simplewallet/data       [no test files]
        simplewallet/config             coverage: 0.0% of statements
        simplewallet/controller         coverage: 0.0% of statements
        simplewallet/service/dao                coverage: 0.0% of statements
        simplewallet/cmd                coverage: 0.0% of statements
ok      simplewallet/controller/validator       (cached)        coverage: 94.7% of statements
        simplewallet/router             coverage: 0.0% of statements
ok      simplewallet/service    (cached)        coverage: 76.2% of statements
ok      simplewallet/util       (cached)        coverage: 56.2% of statements
?       simplewallet/util/errcode       [no test files]
ok      simplewallet/util/db    1.040s  coverage: 32.0% of statements

```

**3. Check goroutine leak**
every TestXXX function will check goroutine leak. Example in code:
```
func TestDeposit(t *testing.T) {
	ConnectDBRedis()
	defer goleak.VerifyNone(t) // check for goroutine leaks
	defer DisconnectDBRedis()
    ........
}

func TestCompareFloat(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks
    ........
}
```

# Questions and Feedback

1. I spend 70% of time on unit test, `validator` can reach 94.7% coverage, but logic `service` can only reach 76.2%. 

2. I test code functions in windows10. Redis installed in ubuntu. I write Dockerfile model for docker, but I have not tested it yet.

3. I don't know how many digits the amount has. So I try `DECIMAL(15, 8)` in mysql. And I check the API request parameter `amount` in `controller/validator` according this.
