## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## start/vercel: start vercel dev server
.PHONY: start/vercel
start/vercel:
	@vercel dev --debug --listen $(APP_PORT)
