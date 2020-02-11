package api

import (
	"fmt"
	"github.com/go-xorm/xorm"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/tv2169145/golang-grpc/api/interceptors"
	"github.com/tv2169145/golang-grpc/api/v1/users"
	pb "github.com/tv2169145/golang-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

func Run(port int, db *xorm.Engine) {
	lst, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcmw.ChainUnaryServer(interceptors.GlobalRepoInjector(db))))
	srv := users.GetRoutes()
	pb.RegisterV1UsersServer(s, srv)
	reflection.Register(s)
	fmt.Printf("server is running on port %d", port)
	if err = s.Serve(lst); err != nil {
		panic(err)
	}
 }
