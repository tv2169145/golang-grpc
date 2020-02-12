package api

import (
	"fmt"
	"github.com/go-xorm/xorm"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/tv2169145/golang-grpc/api/interceptors"
	v1auth "github.com/tv2169145/golang-grpc/api/v1/auth"
	v1user "github.com/tv2169145/golang-grpc/api/v1/users"
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

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmw.ChainUnaryServer(
				interceptors.GlobalRepoInjector(db),
				),
			),
		)
	// 註冊全路由
	initAllRoutes(s)

	reflection.Register(s)
	fmt.Printf("server is running on port %d", port)
	if err = s.Serve(lst); err != nil {
		panic(err)
	}
 }

func initAllRoutes(s *grpc.Server) {
	pb.RegisterV1UsersServer(s, v1user.GetRoutes())
	pb.RegisterV1AuthServer(s, v1auth.GetRoutes())
}
