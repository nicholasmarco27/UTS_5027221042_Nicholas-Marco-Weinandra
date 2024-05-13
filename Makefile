gen:
	@protoc \
		--proto_path=protobuf "protobuf/habittracker.proto" \
		--go_out=common/genproto/habittracker --go_opt=paths=source_relative \
	--go-grpc_out=common/genproto/habittracker --go-grpc_opt=paths=source_relative

run-server:
	@go run cmd/server/main.go $(profile)

run-client:
	@go run cmd/client/main.go $(profile)