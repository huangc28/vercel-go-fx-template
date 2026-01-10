## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## start/vercel: start vercel dev server
.PHONY: start/vercel
start/vercel:
	@vercel dev --debug --listen $(APP_PORT)

## start/inngest: start the inngest dev server
.PHONY: start/inngest
start/inngest:
	PORT=3011 npx inngest-cli@latest dev \
		--no-discovery \
		--poll-interval 10000 \
		-u http://localhost:3010/api/inngest

