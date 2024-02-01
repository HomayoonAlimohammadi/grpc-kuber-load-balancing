PUSH := 0

.PHONY: api-http
api-http:
	cd http && $(MAKE) api	

.PHONY: storage-http
storage-http:
	cd http && $(MAKE) storage	

.PHONY: both-http
both-http: api-http storage-http
	@echo "Made both http!"


.PHONY: api-grpc
api-grpc:
	cd grpc && $(MAKE) api	

.PHONY: storage-grpc
storage-grpc:
	cd grpc && $(MAKE) storage	

.PHONY: both-grpc
both-grpc: api-grpc storage-grpc
	@echo "Made both grpc!"


