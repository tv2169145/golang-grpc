package users

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pb "github.com/tv2169145/golang-grpc/pb"
	repoMocks "github.com/tv2169145/golang-grpc/repos/mocks"
	"github.com/tv2169145/golang-grpc/types"
	"github.com/tv2169145/golang-grpc/utils"
)
var _ = Describe("grpc", func() {
	var (
		globalRepo *repoMocks.MockGlobalRepository
		usersRepo *repoMocks.MockUsersRepo
		mockCtrl *gomock.Controller
		router pb.V1UsersServer
		ctx context.Context
	)
	setupRouter := func() {
		mockCtrl = gomock.NewController(GinkgoT())
		globalRepo = repoMocks.NewMockGlobalRepository(mockCtrl)
		usersRepo = repoMocks.NewMockUsersRepo(mockCtrl)
		router = GetRoutes()
		ctx = utils.SetGlobalRepoOnContext(context.Background(), globalRepo)

		globalRepo.EXPECT().Users().Return(usersRepo).AnyTimes()
	}

	JustBeforeEach(func() {
		setupRouter()
	})

	JustAfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Create", func() {
		It("should return error because of empty request", func() {
			errMsg := "Key: 'CreateUserRequest.NewUser' Error:Field validation for 'NewUser' failed on the 'valid-newUser' tag"
			_, err := router.Create(ctx, &pb.CreateUserRequest{})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should return error because of missing email in request", func() {
			errMsg := `Key: 'CreateUserRequest.email' Error:Field validation for 'email' failed on the 'valid-email' tag
Key: 'CreateUserRequest.firstName' Error:Field validation for 'firstName' failed on the 'valid-firstName' tag
Key: 'CreateUserRequest.lastName' Error:Field validation for 'lastName' failed on the 'valid-lastName' tag
Key: 'CreateUserRequest.password' Error:Field validation for 'password' failed on the 'valid-password' tag
Key: 'CreateUserRequest.confirmPassword' Error:Field validation for 'confirmPassword' failed on the 'valid-confirmPassword' tag`;
			_, err := router.Create(ctx, &pb.CreateUserRequest{
				NewUser: &pb.CreateUser{},
			})
			//spew.Dump(err)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should return error because global repo is missing from context", func() {
			errMsg := "unable to get global repo from context"
			user, err := types.NewUser(&types.TempUser{
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
				Password: "1234",
				ConfirmPassword: "1234",
			})
			Ω(err).To(BeNil())

			_, err = router.Create(context.Background(), &pb.CreateUserRequest{
				NewUser: &pb.CreateUser{
					FirstName: user.FirstName,
					LastName: user.LastName,
					Email: user.Email,
					Password: "1234",
					ConfirmPassword: "1234",
				},
			})
			//spew.Dump(err)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail a database test", func() {
			errMsg := "database is unavailable"

			user, err := types.NewUser(&types.TempUser{
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
				Password: "1234",
				ConfirmPassword: "1234",
			})
			Ω(err).To(BeNil())

			usersRepo.EXPECT().Create(gomock.AssignableToTypeOf(user)).
				Return(errors.New(errMsg)).Times(1).Do(func(*types.User) {
				defer GinkgoRecover()
			})

			_, err = router.Create(ctx, &pb.CreateUserRequest{
				NewUser: &pb.CreateUser{
					FirstName: user.FirstName,
					LastName: user.LastName,
					Email: user.Email,
					Password: "1234",
					ConfirmPassword: "1234",
				},
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should be success", func() {
			user, err := types.NewUser(&types.TempUser{
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
				Password: "1234",
				ConfirmPassword: "1234",
			})
			Ω(err).To(BeNil())

			usersRepo.EXPECT().Create(gomock.AssignableToTypeOf(user)).
				Return(nil).
				Times(1).
				Do(func(*types.User) {
					defer GinkgoRecover()
			})

			res, err := router.Create(ctx, &pb.CreateUserRequest{
				NewUser: &pb.CreateUser{
					FirstName: user.FirstName,
					LastName: user.LastName,
					Email: user.Email,
					Password: "1234",
					ConfirmPassword: "1234",
				},
			})
			Ω(err).To(BeNil())
			Ω(res.GetUser().GetEmail()).To(Equal(user.Email))
		})
	})
})

