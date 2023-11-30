.PHONY: autogen api storage up

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

api: autogen
	@docker build -t post-api-grpc:latest -t homayoonalimohammadi/post-api-grpc:latest .
	@docker push homayoonalimohammadi/post-api-grpc:latest

storage: autogen
	@docker build -t post-storage-grpc:latest -t homayoonalimohammadi/post-storage-grpc:latest ./post_storage
	@docker push homayoonalimohammadi/post-storage-grpc:latest

up: autogen
	@docker compose up