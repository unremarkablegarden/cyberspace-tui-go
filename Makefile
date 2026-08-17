# Export all variables to child processes
.EXPORT_ALL_VARIABLES:
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
APP=cyberspace-cli
OUT=tmp
EXEC_PATH=./cmd/.

.PHONY: build
build:
	@echo -e "$(OK_COLOR)==> Building binary from source$(NO_COLOR)"
	go build -o "$(OUT)/$(APP)" $(EXEC_PATH)
	@echo -e "$(OK_COLOR)==> Done. Check $(OUT)/$(APP)$(NO_COLOR)"

.PHONY: run
run:
	go run $(EXEC_PATH)
