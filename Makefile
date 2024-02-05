PUSH ?= 0

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

.PHONY: install-linkerd
install-linkerd: 
	@if [ $(shell kubectl config current-context) = "kind-kind-test" ]; then \
		linkerd install --crds | kubectl apply -f -; \
		linkerd install | kubectl apply -f -; \
	fi

.PHONY: upgrade-linkerd
upgrade-linkerd: 
	@if [ $(shell kubectl config current-context) = "kind-kind-test" ]; then \
		linkerd upgrade --crds | kubectl apply -f -; \
		linkerd upgrade | kubectl apply -f -; \
	fi

.PHONY: uninstall-linkerd
uninstall-linkerd: 
	@if [ $(shell kubectl config current-context) = "kind-kind-test" ]; then \
		linkerd uninstall | kubectl delete -f -; \
	fi

.PHONY: inject-http-linkerd
inject-http-linkerd: both-http upgrade-linkerd
	@if [ $(shell kubectl config current-context) = "kind-kind-test" ]; then \
		kubectl get deploy post-storage-http post-api-http -o yaml | linkerd inject - | kubectl apply -f -; \
	fi

.PHONY: inject-grpc-linkerd
inject-grpc-linkerd: both-grpc upgrade-linkerd
	@if [ $(shell kubectl config current-context) = "kind-kind-test" ]; then \
		kubectl get deploy post-storage-grpc post-api-grpc -o yaml | linkerd inject - | kubectl apply -f -; \
	fi