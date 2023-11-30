.PHONY: api storage up

api:
	@docker build -t post-api-http:latest -t homayoonalimohammadi/post-api-http:latest .
	@docker push homayoonalimohammadi/post-api-http:latest

storage:
	@docker build -t post-storage-http:latest -t homayoonalimohammadi/post-storage-http:latest ./post_storage
	@docker push homayoonalimohammadi/post-storage-http:latest

up:
	@docker compose up