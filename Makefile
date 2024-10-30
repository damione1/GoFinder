.PHONY: up
up:
	tilt up

.PHONY: down
down:
	tilt down

.PHONY: restart
restart:
	tilt down && tilt up


.PHONY: build-go-proto
build-go-proto:
	rm -f protogen/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --go-grpc_out=pkg/pb --go_out=pkg/pb --proto_path=proto --go-grpc_opt=paths=source_relative \
	--go_opt=paths=source_relative --grpc-gateway_out=pkg/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go-finder \
	./proto/*.proto


.PHONY: generate-sim-key
generate-sim-key:
	TOKEN_SYMMETRIC_KEY=$$(openssl rand -hex 16); \
	echo "TOKEN_SYMMETRIC_KEY=$$TOKEN_SYMMETRIC_KEY" >> .env
