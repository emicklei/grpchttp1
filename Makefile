.PHONY: pb test

pb:
	protofmt -w  io.proto
	protoc io.proto --go_out=. --go-grpc_out=. 