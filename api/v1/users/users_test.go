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

	Describe("FindById", func() {
		It("should fail because the empty request", func() {
			errMsg := "Key: 'FindByIdRequest.id' Error:Field validation for 'id' failed on the 'valid-id' tag"
			_, err := router.FindById(ctx, &pb.FindByIdRequest{})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because get global repo from context error", func() {
			errMsg := "unable to get global repo from context"
			_, err := router.FindById(context.Background(), &pb.FindByIdRequest{
				Id: 1,
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail a database test", func() {
			errMsg := "database is unavailable"

			usersRepo.EXPECT().FindById(int64(1)).
				Return(nil, errors.New(errMsg)).Times(1).Do(func(int64) {
				defer GinkgoRecover()
			})

			_, err := router.FindById(ctx, &pb.FindByIdRequest{
				Id: 1,
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("success", func() {
			user := &types.User{
				ID: 2,
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
				Password: "1234",
			}
			usersRepo.EXPECT().FindById(user.ID).Return(user, nil).Times(1).Do(func(int64){
				defer GinkgoRecover()
			})
			res, err := router.FindById(ctx, &pb.FindByIdRequest{
				Id: user.ID,
			})

			Ω(err).To(BeNil())
			Ω(res.GetUser().GetId()).To(Equal(user.ID))
		})
	})

	Describe("FindByEmail", func() {
		It("should fail because the empty request", func() {
			errMsg := "Key: 'FindByEmailRequest.email' Error:Field validation for 'email' failed on the 'valid-email' tag"
			_, err := router.FindByEmail(ctx, &pb.FindByEmailRequest{})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because get global repo from context error", func() {
			errMsg := "unable to get global repo from context"
			_, err := router.FindByEmail(context.Background(), &pb.FindByEmailRequest{
				Email: "jimmy@gmail.com",
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail a database test", func() {
			errMsg := "database is unavailable"
			usersRepo.EXPECT().FindByEmail("example@gmail.com").
				Return(nil, errors.New(errMsg)).Times(1).Do(func(string) {
				defer GinkgoRecover()
			})
			_, err := router.FindByEmail(ctx, &pb.FindByEmailRequest{
				Email: "example@gmail.com",
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("success", func() {
			user := &types.User{
				ID: 2,
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
				Password: "1234",
			}
			usersRepo.EXPECT().FindByEmail(user.Email).Return(user, nil).Times(1).Do(func(string){
				defer GinkgoRecover()
			})
			res, err := router.FindByEmail(ctx, &pb.FindByEmailRequest{
				Email: user.Email,
			})

			Ω(err).To(BeNil())
			Ω(res.GetUser().GetEmail()).To(Equal(user.Email))
		})
	})

	Describe("Update", func() {
		It("should fail because the empty request", func() {
			errMsg := "Key: 'UpdateUserRequest.id' Error:Field validation for 'id' failed on the 'valid-id' tag"
			_, err := router.Update(ctx, &pb.UpdateUserRequest{})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because get global repo from context error", func() {
			errMsg := "unable to get global repo from context"
			_, err := router.Update(context.Background(), &pb.UpdateUserRequest{
				Id:1,
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because FindById error", func() {
			errMsg := "database error"
			usersRepo.EXPECT().FindById(int64(1)).Return(nil, errors.New(errMsg)).Do(func(int64) {
				defer GinkgoRecover()
			})
			_, err := router.Update(ctx, &pb.UpdateUserRequest{
				Id:1,
				FirstName: "linn",
				LastName: "Jimmy",
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because database error", func() {
			errMsg := "set password error"
			user := &types.User{
				ID: 1,
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
			}
			usersRepo.EXPECT().FindById(int64(1)).Return(user, nil).Times(1).Do(func(int64) {
				defer GinkgoRecover()
			})

			usersRepo.EXPECT().Update(user).Return(errors.New(errMsg)).Do(func(*types.User) {
				defer GinkgoRecover()
			})
			_, err := router.Update(ctx, &pb.UpdateUserRequest{
				Id: 1,
				FirstName: "Linn",
				LastName: "jimmy",
				NewPassword: "1234",
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("success", func() {
			user := &types.User{
				ID: 1,
				FirstName: "Linn",
				LastName: "jimmy",
				Email: "jimmy@gmail.com",
			}
			usersRepo.EXPECT().FindById(int64(1)).Return(user, nil).Times(1).Do(func(int64) {
				defer GinkgoRecover()
			})

			usersRepo.EXPECT().Update(user).Return(nil).Do(func(*types.User) {
				defer GinkgoRecover()
			})

			res, err := router.Update(ctx, &pb.UpdateUserRequest{
				Id: user.ID,
				FirstName: user.FirstName,
				LastName: user.LastName,
				NewPassword: "1234",
			})
			Ω(err).To(BeNil())
			Ω(res).NotTo(BeNil())
			Ω(res.GetUser().GetId()).To(Equal(user.ID))
		})
	})
})

