# Contributing

Clone the repository:

```shell
git clone git@github.com:gruyaume/lesvieux.git
```

Run Go unit tests:

```shell
go test ./...
```

Generate the sqlc code:

```shell
sqlc generate
```

Build the frontend:

```shell
npm install --prefix ui
npm run build --prefix ui
```

View the frontend:

```shell
go run cmd/lesvieux/main.go -config lesvieux.yaml
```

Navigate to https://localhost:8000.
