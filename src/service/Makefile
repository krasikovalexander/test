.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/list functions/list/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/compare functions/compare/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/compare-routes functions/compare-routes/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/rank functions/rank/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
