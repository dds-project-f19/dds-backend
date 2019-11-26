# **Distributed and Decentralised Systems** Course Project
[![Go Report Card](https://goreportcard.com/badge/github.com/dds-project-f19/dds-backend)](https://goreportcard.com/report/github.com/dds-project-f19/dds-backend)
[![MIT](https://img.shields.io/github/license/dds-project-f19/dds-backend)](https://raw.githubusercontent.com/dds-project-f19/dds-backend/master/LICENSE)
![SIZE](https://img.shields.io/github/repo-size/dds-project-f19/dds-backend)

### Build and Run:

#### 1 Configure database link (mysql) /config/config.yaml:

#### 2 Get dependencies

```shell script
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jinzhu/gorm
go get -u github.com/gin-gonic/gin
go get -u github.com/gin-contrib/cors
go get -u github.com/gin-contrib/static
go get -u github.com/go-telegram-bot-api/telegram-bot-api

```
#### 3 Setup Environment Variables
**DDS_TELEGRAM_BOT_APIKEY** - api key for telegram bot
#### 4 Get server

INSIDE $GOPATH/src
```shell script
git clone https://github.com/dds-project-f19/dds-backend
```

#### 5 Build

```shell script
go install $GOPATH/src/dds-backend
```

#### 6 Run

Now you can launch with executable located at `$GOPATH/bin/` called `dds-backend` (.exe for Windows)
```shell script
$ dds-backend --help
Usage of dds-backend:
  -dbaddress string
         (default "127.0.0.1")
  -dblogin string
         (default "root")
  -dbname string
         (default "ddstest")
  -dbpassword string
         (default "ddspassword14882")
  -dbport string
         (default "3306")
```


### API Description (to be moved to wiki later):

```golang
// POST /common/login
// HEADERS: {}
// {"username":"123", "password":"456"}
// 200: {"token":"1234567"}
// 400,403: {"message":"123"}

// GET /common/telegram_join_link
// HEADERS: {Authorization: token}
// {}
// 200: {"link":"t.me/bot_link/start=regkey123"}
// 401, 500: {"message":"123"}

// GET /worker/get
// HEADERS: {Authorization: token}
// {}
// 200: {"username":"required", "name":"", "surname":"", "phone":"", "address":""}
// 401,404: {"message":"123"}

// PATCH /worker/update
// HEADERS: {Authorization: token}
// {"username":"required", "name":"", "surname":"", "phone":"", "address":""}
// 200: {}
// 400,401,404: {"message":"123"}

// POST /worker/take_item
// HEADERS: {Authorization: token}
// {"itemtype":"123", "slot":"123"}
// 201: {"message":"request done, blah blah"}
// 400,401,500: {"message":"123"}

// POST /worker/return_item
// HEADERS: {Authorization: token}
// {"itemtype":"123", "slot":"123"}
// 201: {"message":"request done, blah blah"}
// 400,401,500: {"message":"123"}

// GET /worker/list_available_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"itemtype":"123","count":77}]}
// 401,500: {"message":"123"}

// GET /worker/list_taken_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"takenby":"username","itemtype":"123","assignedtoslot":"123"}]}
// 401,500: {"message":"123"}

// POST /manager/register_worker
// {"username":"required", "password":"required", "name":"", "surname":"", "phone":"", "address":""}
// 201: {"token":"1234567"}
// 400,409,500: {"message":"123"}

// GET /manager/list_workers
// HEADERS: {Authorization: token}
// {}
// 200: {"users":[{"username":""...}]}
// 401,500: {"message":"123"}

// DELETE /manager/remove_worker/{username}
// HEADERS: {Authorization: token}
// {}
// 200: {}
// 400,401,404,500: {"message":"123"}

// PATCH /manager/add_available_items
// HEADERS: {Authorization: token}
// {"itemtype":"123","count":77}
// 200: {}
// 400,401,500: {"message":"123"}

// PATCH /manager/remove_available_items
// HEADERS: {Authorization: token}
// {"itemtype":"123","count":77}
// 200: {}
// 400,401,500: {"message":"123"}

// GET /manager/list_available_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"itemtype":"123","count":77}]}
// 401,500: {"message":"123"}

// GET /manager/list_taken_items
// HEADERS: {Authorization: token}
// {}
// 200: {"items":[{"takenby":"username","itemtype":"123","assignedtoslot":"123"}]}
// 401,500: {"message":"123"}

// POST /admin/register_manager
// HEADERS: {}
// {"username":"required", "password":"required", "gametype":"required", "name":"", "surname":"", "phone":"", "address":""}
// 201: {}
// 400,401,409,500: {"message":"123"}
```
