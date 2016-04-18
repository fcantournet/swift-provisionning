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
	scp ${project} s-osncls-0003.adm.stg0.aub.cloudwatt.net:~/
	scp ${project} s-osncls-0004.adm.stg0.aub.cloudwatt.net:~/
	scp ${project} s-osncls-0003.adm.stg1.aub.cloudwatt.net:~/
	scp ${project} s-osncls-0004.adm.stg1.aub.cloudwatt.net:~/
