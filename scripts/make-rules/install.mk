

.PHONY: install.verify.%
install.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) install.$*; fi

.PHONY: install.reshim
install.reshim:
	@if command -v asdf &>/dev/null; then asdf reshim golang; fi

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

.PHONY: install.gsemver
install.gsemver:
	@$(GO) install github.com/arnaud-deprez/gsemver@latest

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: install.github-release
install.github-release:
	@$(GO) install github.com/github-release/github-release@latest

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
	@$(MAKE) install.reshim
	@golangci-lint completion $$(basename $$SHELL) > $(HOME)/.golangci-lint.$$(basename $$SHELL)
	@if ! grep -q .golangci-lint.$$(basename $$SHELL) $(HOME)/.$$(basename $$SHELL)rc; then echo "source \$$HOME/.golangci-lint.$$(basename $$SHELL)" >> $(HOME)/.$$(basename $$SHELL)rc; fi