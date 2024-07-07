### dependencies

- [serverless](https://www.npmjs.com/package/serverless)
- [serverless offline](https://www.npmjs.com/package/serverless-offline)
- [AWS Lambda for Go](https://github.com/aws/aws-lambda-go)


### install dependencies

```
npm install -g serverless
npm install
go mod tidy
docker pull public.ecr.aws/lambda/go
```

### build

```
make build
```

### run service
```
make start
```

### simple ping request

```
 curl http://0.0.0.0:3000/ping
```
