PUSH ?= 0

autogen: 
	@protoc \
	--proto_path=post_storage/proto \
	--go_out=post_storage \
	--go_opt=Mpost_storage.proto=proto/autogen/post_storage \
	--go-grpc_opt=Mpost_storage.proto=proto/autogen/post_storage \
	--go-grpc_out=post_storage \
	post_storage.proto

	@protoc \
	--proto_path=post_api/proto \
	--go_out=post_api \
	--go_opt=Mpost_api.proto=proto/autogen/post_api \
	--go-grpc_opt=Mpost_api.proto=proto/autogen/post_api \
	--go-grpc_out=post_api \
	post_api.proto

	@echo "Generated files from proto successfully!"

.PHONY: api
api: autogen
	@docker build -t post-api-grpc:latest -t homayoonalimohammadi/post-api-grpc:latest .
	@if [ "$$PUSH" = "1" ]; then \
		docker push homayoonalimohammadi/post-api-grpc:latest; \
	fi
	kind load docker-image post-api-grpc:latest


.PHONY: storage
storage: autogen
	@docker build -t post-storage-grpc:latest -t homayoonalimohammadi/post-storage-grpc:latest ./post_storage
	@if [ "$$PUSH" = "1" ]; then \
		docker push homayoonalimohammadi/post-storage-grpc:latest; \
	fi
	kind load docker-image post-storage-grpc:latest

.PHONY: both
both: api storage
	@echo "Made both!"

.PHONY: up
up: autogen
	@docker compose up