![alt tag](https://upload.wikimedia.org/wikipedia/commons/2/23/Golang.png)

## POC for Go Lab - 2023

**This is auditt_98's fork of Golang Gin Boilerplate**

Bootstrapping projects with built-in batteries included, such as:
  - Authentication (OIDC with Zitadel)
	- Query Helpers
	- CockroachDB
	- Presigned URLS
	- Validators
	- Middlewares (CORS, RequestID, etc)
	- Env support
	- SSL support

### Installation

```
$ go get
$ go install
```


### Environment variables

Multiple environment variables are supported (dev, staging, prod)
There are a few ways to load the default env

#### 1. Using the `ONTHEGO_ENV` environment variable 
```
	export ONTHEGO_ENV=dev
```

#### 2. Using the env=[target] flag
```
	air env=dev
```

#### 3. Default to dev
```
	air
```

* Note that env=[target] flag has the highest priority, followed by ONTHEGO_ENV, and then defaulting to dev

#### Replace the following in the `.env` file
```
ENV=LOCAL
PORT=9000
SSL=FALSE
API_VERSION=1.0
DB_DRIVER=postgres
DB_USER="root"
DB_PASS="postgres"
DB_HOST="localhost"
DB_PORT="26257"
DB_NAME="onthego"
REDIS_SECRET=""
REDIS_HOST=127.0.0.1:6379
REDIS_PASSWORD=
ZITADEL_USER_EMAIL=""
ZITADEL_USERNAME="core_human_user"
ZITADEL_PASSWORD="CoreHumanPassword1!"
ZITADEL_DOMAIN=http://localhost:8080
API_DOMAIN=http://host.docker.internal:9000
FILE_UPLOAD_PATH=uploads
UPLOAD_DRIVER=local
SIGNED_URL_SECRET=CatMeowMeow123
ZITADEL_DATABASE_COCKROACH_HOST=crdb
ZITADEL_EXTERNALSECURE=false
ZITADEL_LOG_LEVEL=info
ZITADEL_DEFAULTINSTANCE_ORG_NAME=OnTheGo
ZITADEL_DEFAULTINSTANCE_DOMAINPOLICY_SMTPSENDERADDRESSMATCHESINSTANCEDOMAIN=false
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_SMTP_HOST=
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_SMTP_USER=
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_SMTP_PASSWORD=
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_TLS=true
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_FROM=
ZITADEL_DEFAULTINSTANCE_SMTPCONFIGURATION_FROMNAME=
```

#### Currently, the supported drivers are:
- mysql
- postgres (support Postgres and CockroachDB) (Recommended)
- mssql

## Building Your Application

```
$ go build -v
```

```
$ ./gin-boilerplate
```

## Testing Your Application

```
$ go test -v ./tests/*
```

## Contribution

You are welcome to contribute to keep it up to date and always improving!


## Credit

This project is forked from [Golang Gin Boilerplate]()



## License

(The MIT License)

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
