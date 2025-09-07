.PHONY: build clean install uninstall

LOCAL_BIN = /usr/local/bin/

build:
	go build -o jenkinsfile-validator main.go

clean:
	@echo "Cleaning jenkinsfile-validator from local repo"
	rm -rf jenkinsfile-validator

install:
	@echo "Installing jenkinsfile validator"
	go build -o jenkinsfile-validator main.go
	sudo mv jenkinsfile-validator $(LOCAL_BIN)

clean-config:
	sudo rm -rf ~/.validator_config.json 
uninstall:
	sudo rm -rf $(LOCAL_BIN)jenkinsfile-validator
