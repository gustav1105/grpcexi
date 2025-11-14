build_proto:
	protoc --go_out=protos --go-grpc_out=protos protos/contact.proto
	
