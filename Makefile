test:
	go clean -testcache
	go test ./...

vtest:
	go clean -testcache
	go test -v ./...

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

server:
	go run cmd/main.go -- start

debug:
	dlv debug -- start

migrate:
	go run cmd/main.go migrate

resetdb:
	go run cmd/main.go rollback `/bin/ls -1 ./internal/repository/migrations/*down.sql | wc -l`
	go run cmd/main.go migrate

seed:
	go run cmd/main.go seed

todo:
	grep -Rin --include="*go" "TODO" * 