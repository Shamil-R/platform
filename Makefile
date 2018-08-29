PROJECTNAME=$(shell basename "$(PWD)")

PROJECTPATH=$(shell pwd)

CMDPATH=$(PROJECTPATH)/platform
CMDMAIN=$(CMDPATH)/main.go

go-generate:
	@echo "  >  Generate..."
	go generate ../...

go-run:
	@echo "  >  Run project..."
	go run $(CMDMAIN) codegen