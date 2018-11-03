BINDIR=bin
BIN_WORKER=worker
BIN_SANDBOX=sandbox

binary: clean
	if [ -z "$$GOPATH" ]; then \
	  echo "GOPATH is not set"; \
	  exit 1; \
	fi

	@mkdir -p ./$(BINDIR)
	GOOS=linux GOARCH=amd64 go build -v \
	    -o $(BINDIR)/$(BIN_WORKER) ./main/worker

	GOOS=linux GOARCH=amd64 go build \
		-o $(BINDIR)/$(BIN_SANDBOX) ./main/sandbox

clean:
	rm -rf "$(BINDIR)"

.PHONY: clean binary
