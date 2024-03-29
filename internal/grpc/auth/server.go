package auth

import (
	ssov1 "github.com/igorakimy/grpc-sso-auth-service/gen/protobuf/v1/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	ssov1.UnimplementedAuthServer
	authService AuthenticationService
}

func RegisterServer(grpcServer *grpc.Server, authService AuthenticationService) {
	reflection.Register(grpcServer)
	ssov1.RegisterAuthServer(grpcServer, &server{
		authService: authService,
	})
}
