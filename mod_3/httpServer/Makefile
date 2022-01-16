export tag=v1.2
root:
	export ROOT=github.com/jzhang026/cloud-native-stack/tree/main/mod_3

release:
	echo "building httpserver container"
	docker build -t msybe/httpserver:${tag} .

push: release
	echo "pushing msybe/httpserver"
	docker push msybe/httpserver:${tag}
