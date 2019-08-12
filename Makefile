# Go parameters

all: init

build:
	go build qms.mgmt.api

clean:
	go clean qms.mgmt.api

init:export http_proxy=http://127.0.0.1:7777/pac
init:export https_proxy=http://127.0.0.1:7777/pac	
init:
	dep ensure
	
update:export http_proxy=http://127.0.0.1:7777/pac
update:export https_proxy=http://127.0.0.1:7777/pac	
update:
	dep ensure -update