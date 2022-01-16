export tag=v1.1
root:
	export ROOT=github.com/jzhang026/cloud-native-stack/tree/main/mod_3

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t msybe/httpserver:${tag} .

push: release
	echo "pushing msybe/httpserver"
	docker push msybe/httpserver:${tag}
