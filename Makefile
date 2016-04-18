project=swift-provisionning
version=$(shell git describe --tags)

all: ${project}

${project}: build

build:
	go build -o ${project}

static:
	./build-static

build-indocker:
	docker run --rm --name=${project}-build -v $(shell pwd)/:/build/code ${builddockerimage} make static

dockerimage: build-indocker
	docker build -t ${rundockerimage}:${version} .
	docker push ${rundockerimage}

deploy: build
	scp ${project} p-osncls-0005.adm.prd2.rue.cloudwatt.net:~/
