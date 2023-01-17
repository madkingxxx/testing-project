CURRENT_DIR=${PWD}

gen-proto:
	protoc -I=${CURRENT_DIR}/proto --go_out=${CURRENT_DIR} \
		--go-grpc_out=${CURRENT_DIR} ${CURRENT_DIR}/proto/*.proto

gen-swag:
	swag init -g internal/delivery/http/api.go -o docs
.PHONY: gen-proto gen-swag