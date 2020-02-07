package users

import (
	"context"
	pb "github.com/tv2169145/golang-grpc/pb"
	"github.com/tv2169145/golang-grpc/types"
	"github.com/tv2169145/golang-grpc/utils"
)

type grpcHandler struct {

}

func GetRoutes() pb.V1UsersServer {
	return &grpcHandler{}
}

/*
	rpc Create(CreateUserRequest) returns(UserReply) {}
    rpc FindById(FindByIdRequest) returns(UserReply) {}
    rpc FindByEmail(FindByEmailRequest) returns(UserReply) {}
    rpc Update(UpdateUserRequest) returns(UserReply) {}
 */

func(h *grpcHandler) Create(ctx context.Context, req *pb.CreateUserRequest) (res *pb.UserReply, err error) {
	res = new(pb.UserReply)

	if err = pb.Validate(req); err != nil {
		return
	}
	globalRepo, err := utils.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return
	}
	newUser, err := types.NewUser(&types.TempUser{
		FirstName: req.GetNewUser().GetFirstName(),
		LastName: req.GetNewUser().GetLastName(),
		Email: req.GetNewUser().GetEmail(),
		Password: req.GetNewUser().GetPassword(),
		ConfirmPassword: req.GetNewUser().GetConfirmPassword(),
	})

	if err = globalRepo.Users().Create(newUser); err != nil {
		return
	}

	res.User = newUser.ToProtoBuf()
	return
}

func(h *grpcHandler) FindById(ctx context.Context, req *pb.FindByIdRequest) (res *pb.UserReply, err error) {
	res = new(pb.UserReply)
	return
}

func(h *grpcHandler) FindByEmail(ctx context.Context, req *pb.FindByEmailRequest) (res *pb.UserReply, err error) {
	res = new(pb.UserReply)
	return
}

func(h *grpcHandler) Update(ctx context.Context, req *pb.UpdateUserRequest) (res *pb.UserReply, err error) {
	res = new(pb.UserReply)
	return
}

