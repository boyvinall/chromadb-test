.PHONY: all build start stop clean lint

all: build lint

define PROMPT
	@echo
	@echo "**********************************************************"
	@echo "*"
	@echo "*   $(1)"
	@echo "*"
	@echo "**********************************************************"
	@echo
endef

build:
	$(call PROMPT, $@)
	go build -o chromadb-test .

start:
	$(call PROMPT, $@)
	docker-compose up -d

stop:
	$(call PROMPT, $@)
	docker-compose down

clean:
	$(call PROMPT, $@)
	rm -f chromadb-test

lint:
	$(call PROMPT, $@)
	golangci-lint run

