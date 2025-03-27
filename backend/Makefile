.PHONY: swag ### generate swagger docs
swag:
	go tool swag init -g 'internal/app/app.go' --parseInternal --parseDependency

.PHONY: swag-fmt ### fmt swagger annotations
swag-fmt:
	go tool swag fmt
