user	:=	$(shell whoami)
rev 	:= 	$(shell git rev-parse --short HEAD)

all:
	@cd cmd/protoc-gen-go-http && go build && cd - &> /dev/null


.PHONY: install
install:
	@cp ./cmd/protoc-gen-go-http/protoc-gen-go-http /usr/local/bin
	@echo "install finished"

.PHONY: uninstall
uninstall:
	$(shell for i in `which -a protoc-gen-validate | grep -v '/usr/bin/protoc-gen-go-http' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	@echo "uninstall finished"

test:
	@echo $(user) $(rev)

dependency:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest