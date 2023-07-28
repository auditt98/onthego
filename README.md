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
- Enviroment support
- Unit test
- And few other important utilties to kickstart any project
- Drivers for a few databases built-in (Gorm drivers)
- Built-in high level Query Engine on top of Gorm that provides selecting syntax similar to MongoDB


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
- postgres
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

### Running Your Application

This project supports Live Reloading with [Air](https://github.com/cosmtrek/air)

To start the development of the project, run the following command:

```
  air
```


Generate SSL certificates (Optional)

> If you don't SSL now, change `SSL=TRUE` to `SSL=FALSE` in the `.env` file

```
$ mkdir cert/
```

```
$ sh generate-certificate.sh
```

> Make sure to change the values in .env for your databases

```
$ go run *.go
```

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
