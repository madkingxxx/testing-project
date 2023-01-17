CURRENT_DIR=${PWD}

gen-proto:
	protoc -I=${CURRENT_DIR}/proto --go_out=${CURRENT_DIR} \
		--go-grpc_out=${CURRENT_DIR} ${CURRENT_DIR}/proto/*.proto

.PHONY: gen-proto