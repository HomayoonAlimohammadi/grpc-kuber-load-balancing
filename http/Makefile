PUSH ?= 0

.PHONY: api
api:
	docker build -t post-api-http:latest -t homayoonalimohammadi/post-api-http:latest .
	@if [ "$$PUSH" = "1" ]; then \
		docker push homayoonalimohammadi/post-api-http:latest; \
	fi
	kind load docker-image post-api-http:latest

.PHONY: storage
storage:
	docker build -t post-storage-http:latest -t homayoonalimohammadi/post-storage-http:latest ./post_storage
	@if [ "$$PUSH" = "1" ]; then \
		docker push homayoonalimohammadi/post-storage-http:latest; \
	fi
	kind load docker-image post-storage-http:latest

.PHONY: both
both: api storage 
	@echo "Made both!"

.PHONY: up
up:
	@docker compose up