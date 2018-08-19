PROJECTNAME=$(shell basename "$(PWD)")

PROJECTPATH=$(shell pwd)

CMDPATH=$(PROJECTPATH)/platform
CMDMAIN=$(CMDPATH)/main.go

go-run:
	@echo "  >  Run project..."
	@cd $(CMDPATH) && go run ./main.go codegen