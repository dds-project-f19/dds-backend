# **Distributed and Decentralised Systems** Course Project
[![Go Report Card](https://goreportcard.com/badge/github.com/dds-project-f19/dds-backend)](https://goreportcard.com/report/github.com/dds-project-f19/dds-backend)
[![MIT](https://img.shields.io/github/license/dds-project-f19/dds-backend)](https://raw.githubusercontent.com/dds-project-f19/dds-backend/master/LICENSE)
![SIZE](https://img.shields.io/github/repo-size/dds-project-f19/dds-backend)
### Configure database link (mysql) /config/config.yaml:
```
dsn: "login:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
```

### Build and Run:
```shell script
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jinzhu/gorm
go get -u github.com/dds-project-f19/dds-backend
```

all packages should be now located at `$GOPATH/src`

```shell script
go install your/go/path/src/dds-backend
```

now you can launch with executable located at `$GOPATH/bin/` called `dds-backend` (.exe for Windows)


### API Description (to be moved to wiki later):
TODO:
* Add auth tokens
* Decide on auth model
* Add admin user

Ping:
> GET /ping

Expected:
```
PONG
```

Register User:
> POST /users/register

Request:
```json
{
	"username": "username",
	"name": "name",
	"password": "password"
}
```
Expected:
```
{
    "message": "User created successfully",
    "status": "success"
}
```

Get Users list:
> GET /users/list

Expected:
```json
{
    "data": [
        {
            "id": 1,
            "created_at": "2019-09-24T20:49:00+03:00",
            "updated_at": "2019-09-24T20:49:00+03:00",
            "username": "myusername",
            "name": "myname"
        },
        {
            "id": 2,
            "created_at": "2019-09-24T20:55:01+03:00",
            "updated_at": "2019-09-24T20:56:35+03:00",
            "username": "myusername2",
            "name": "myname2"
        }
    ],
    "status": "success"
}
```

Delete User:
> DELETE /users/remove/123

Update User details:
> PATCH /users/edit/123
