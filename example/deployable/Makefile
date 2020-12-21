CODE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CWD := ${CURDIR}

SRC=$(CODE_DIR)/src/*.go $(CODE_DIR)/../../*.go

FUNCTIONS = \
	$(CODE_DIR)/build/handler

all: deploy

$(CODE_DIR)/build/handler: $(SRC)
	cd $(CODE_DIR) && env GOARCH=amd64 GOOS=linux go build -o $(CODE_DIR)/build/handler $(CODE_DIR)/src/ && cd $(CWD)

build: $(FUNCTIONS)

deploy: $(FUNCTIONS)
	cdk deploy
