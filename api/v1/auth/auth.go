package auth

import (
	"context"
	pb "github.com/tv2169145/golang-grpc/pb"
	"github.com/tv2169145/golang-grpc/utils"
)

type grpcHandler struct {

}

func InitRoutes() *grpcHandler {
	return &grpcHandler{}
}

func(h *grpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (res *pb.LoginReply, err error) {
	res = new(pb.LoginReply)
	if err = pb.Validate(req); err != nil {
		return
	}

	globalRepo, err := utils.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return
	}
	usersRepo := globalRepo.Users()
	authRepo := globalRepo.Auth()
	// 取出user
	user, err := usersRepo.FindByEmail(req.GetEmail())
	if err != nil {
		return
	}
	// 檢查password and visible
	if err = user.Authenticate(req.GetPassword()); err != nil {
		return
	}
	claims := authRepo.GetNewClaims(user.Email, map[string]interface{}{
		"user": user,
	})

	token, err := authRepo.GetSignedToken(claims)
	if err != nil {
		return
	}

	res.Token = token
	return
}
