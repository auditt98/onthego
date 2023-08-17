![alt tag](https://upload.wikimedia.org/wikipedia/commons/2/23/Golang.png)

## Golang Gin

**This is auditt_98's fork of Golang Gin Boilerplate**

It provides a quick way to bootstrap your project with a predefined structure with support for multiple databases, Redis, OIDC authentication
(to be implemented), and more.

### Configured with

- [gorm](https://gorm.io/): Gorm
- [jwt-go](https://github.com/golang-jwt/jwt): JSON Web Tokens (JWT) as middleware
- [go-redis](https://github.com/go-redis/redis): Redis support for Go
- Go Modules
- Built-in **Custom Validators**
- Built-in **CORS Middleware**
- Built-in **RequestID Middleware**
- SSL Support
- Environment support
- Unit test
- And few other important utilties to kickstart any project
- Drivers for a few databases built-in (Gorm drivers)
- High level Query Engine on top of Gorm that provides selecting syntax similar to MongoDB (WIP)
- Code Scaffolding (WIP)
- Batteries included with [https://docs.rkdev.info/](RK-Boot)


### On my TODO list
- Finish implementing the Query Engine to support Projection, Population
- Implementing OIDC authentication (Using Zitadel)
- Dockerize

### Installation

```
$ go get onthego
```

```
$ cd $GOPATH/src/onthego
```

```
$ go mod init
```

```
$ go install
```

### Database Driver
You can select the database driver you want to use by changing the `DB_DRIVER` in the `.env` file


#### Currently, the supported drivers are:
- mysql
- postgres (support Postgres and CockroachDB)
- mssql

#### Replace the following in the `.env` file
```
DB_DRIVER=
DB_USER=
DB_PASS=
DB_HOST=
DB_PORT=
DB_NAME=
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

### Authentication
This project aims to use Zitadel as the authentication provider. (WIP)
Zitadel will requires an instance of Zitadel and a single-node CockroachDB instance to run. (docker-compose.yaml)

Default credentials for Zitadel:
```
	username: zitadel-admin@zitadel.localhost
	password: Password1!
```

#### Setup Authorization



### Running Your Application

- Step 1: Run docker compose up to start the database and Zitadel
```
	docker-compose up
```

- Step 2: Setup Zitadel


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

## Using the Query Engine
```
params := db.QueryParams{
	Where: db.WhereParams{
		Or: []db.WhereParams{
			{
				Attr: map[string]db.AttributeParams{
					"id": {
						Eq: 2,
					},
				},
			},
		},
	},
	OrderBy: []string{"-id"},
	Limit:   2,
	Offset:  0,
	Populate: []string{
		"Articles",
	},
	Picks: []string{
		"id",
		"email",
	},
}
var users []models.User
db.Query("users", params, &users)
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
