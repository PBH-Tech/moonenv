.PHONY: build run clean

build:
	$(call buildFile,"endpoints/upload-file") \
	&& $(call buildFile,"endpoints/download-file") \
	&& $(call buildFile,"endpoints/orchestrator")

run: build
	sls deploy --stage dev

clean:
	sls remove --stage dev

define buildFile
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .bin/${1}/bootstrap ${1}/main.go
endef