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
	scp ${project} p-osncls-0000.adm.prd2.rue.cloudwatt.net:~/
	scp ${project} p-osncls-0001.adm.prd2.rue.cloudwatt.net:~/
	scp ${project} p-osncls-0002.adm.prd2.rue.cloudwatt.net:~/

deploy-int5: build
	scp ${project} i-osncls-0000.adm.int5.aub.cloudwatt.net:~/
	scp ${project} i-osncls-0001.adm.int5.aub.cloudwatt.net:~/
	scp ${project} i-osncls-0002.adm.int5.aub.cloudwatt.net:~/
