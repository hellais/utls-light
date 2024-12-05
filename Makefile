sync-upstream:
	scripts/sync-upstream.sh

test:
	go test

.PHONY: test sync-upstream
