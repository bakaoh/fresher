
gen:
	protoc -I/usr/local/include -I. \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types:. \
		CounterService.proto

bin:
	go test -c

run:
	PORT=56789 ./counter.test -test.cpu 4 -test.bench '.*'