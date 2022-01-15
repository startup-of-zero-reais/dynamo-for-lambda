GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

usage:
	@printf "\nComandos disponíveis:\n"
	@printf "\tmock\t- Gera os arquivos de mock\n"
	@printf "\ttest\t- Roda os testes com coverage\n"

.PHONY: mock
mock:
	./local-scripts/install-mockery
	source ~/.bashrc
	mockery --all --keeptree

.PHONY: test
test:
	ENVIRONMENT=testing $(GOTEST) -coverprofile=./tests/coverage-report.out ./...
	$(GOCOVER) -func=./tests/coverage-report.out
	$(GOCOVER) -html=./tests/coverage-report.out