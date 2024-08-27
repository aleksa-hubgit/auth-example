protoc: clean
	cd proto && protoc --go_out=../protogen/golang \
	--go_opt=paths=source_relative --go-grpc_out=../protogen/golang\
	 --go-grpc_opt=paths=source_relative -I=$(PWD) proto/hello/helloworld.proto
clean:
	rm -rf protogen
