## Main goals:
- Support the creation of Handlers, Models and Services through multiple interfaces: CLI, GUI, REST API, etc.

- Use Air to implements hot reload

- Do I need to keep track of the previous tree?

```
go run ./cli/cli.go gen handler
```

### Generate API Version
- Add new folder to './handlers/' with the name of the version
- Write a new group to main.go
```
  v1: r.Group(prefix) {

  }
```


### Generate API Group

### Generate API Handler

### Generate API Action