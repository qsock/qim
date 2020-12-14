LINUXARG=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

BUILDARG=-ldflags " -s -X main.buildtime=`date '+%Y-%m-%d_%H:%M:%S'` -X main.githash=`git rev-parse HEAD`"

p:
	protoc -I ./proto --gogofaster_out=plugins=grpc:. errmsg.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. model.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. stream.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. ret.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. event.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. passport.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. id.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. user.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. file.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. msg.proto;
	protoc -I ./proto --gogofaster_out=plugins=grpc:. ws.proto;
	cp -fR ./github.com/qsock/qim/lib/* ./lib/;
	rm -rf ./qsock.com/;

method:
	python3 tools/method_gen.py `pwd`;go fmt `pwd`'/lib/method';
	python3 tools/proto_gen.py `pwd`;go fmt `pwd`'/api_gateway/controller';

swag:
	cd api_gateway;swag init --parseDependency true;cd -;

api:
	go install ${BUILDARG} github.com/qsock/qim/api_gateway;

passport_server:
	go install ${BUILDARG} github.com/qsock/qim/server/passport_server;

id_server:
	go install ${BUILDARG} github.com/qsock/qim/server/id_server;

user_server:
	go install ${BUILDARG} github.com/qsock/qim/server/user_server;

file_server:
	go install ${BUILDARG} github.com/qsock/qim/server/file_server;

msg_server:
	go install ${BUILDARG} github.com/qsock/qim/server/msg_server;

ws_server:
	go install ${BUILDARG} github.com/qsock/qim/server/ws_server;
