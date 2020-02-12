package auth

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pascaldekloe/jwt"
	pb "github.com/tv2169145/golang-grpc/pb"
	repoMocks "github.com/tv2169145/golang-grpc/repos/mocks"
	"github.com/tv2169145/golang-grpc/types"
	"github.com/tv2169145/golang-grpc/utils"
)

var _ = Describe("", func() {
	var (
		globalRepo *repoMocks.MockGlobalRepository
		usersRepo *repoMocks.MockUsersRepo
		authRepo *repoMocks.MockAuthRepo
		mockCtrl *gomock.Controller
		router pb.V1AuthServer
		ctx context.Context

		validRequest *pb.LoginRequest
		validUser *types.User
		validClaims *jwt.Claims
	)
	setupRouter := func() {
		mockCtrl = gomock.NewController(GinkgoT())
		globalRepo = repoMocks.NewMockGlobalRepository(mockCtrl)
		usersRepo = repoMocks.NewMockUsersRepo(mockCtrl)
		authRepo = repoMocks.NewMockAuthRepo(mockCtrl)
		router = GetRoutes()
		ctx = utils.SetGlobalRepoOnContext(context.Background(), globalRepo)

		globalRepo.EXPECT().Auth().Return(authRepo).AnyTimes()
		globalRepo.EXPECT().Users().Return(usersRepo).AnyTimes()
	}

	setupData := func() {
		validRequest = &pb.LoginRequest{
			Email: "jimmy@gmail.com",
			Password: "1234",
		}

		validUser = &types.User{
			ID: 1,
			Email: "jimmy@gmail.com",
			FirstName: "Linn",
			LastName: "jimmy",
			Visible: true,
		}
		validUser.SetPassword(validRequest.GetPassword())

		validClaims = &jwt.Claims{
			Set: map[string]interface{}{
				"user": validUser,
			},
			Registered: jwt.Registered{
				Subject: validUser.Email,
			},
		}
	}

	JustBeforeEach(func() {
		setupRouter()
		setupData()
	})
	JustAfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Login", func() {
		It("should fail because email is empty", func() {
			errMsg := "Key: 'LoginRequest.email' Error:Field validation for 'email' failed on the 'valid-email' tag"
			_, err := router.Login(ctx, &pb.LoginRequest{
				Email: "",
				Password: "1234",
			})
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because get global repo from context missing", func() {
			errMsg := "unable to get global repo from context"
			_, err := router.Login(context.Background(), validRequest)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because database error", func() {
			errMsg := "database is unavailable"
			usersRepo.EXPECT().FindByEmail(string(validRequest.GetEmail())).Return(nil, errors.New(errMsg)).Times(1).Do(func(string) {
				defer GinkgoRecover()
			})

			_, err := router.Login(ctx, validRequest)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because password not match", func() {
			errMsg := "crypto/bcrypt: hashedPassword is not the hash of the given password"
			validRequest.Password = "123"

			usersRepo.EXPECT().FindByEmail(validRequest.GetEmail()).Return(validUser, nil).Times(1).Do(func(string) {
				defer GinkgoRecover()
			})

			_, err := router.Login(ctx, validRequest)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("should fail because can not get token", func() {
			errMsg := "asdasd"
			usersRepo.EXPECT().FindByEmail(validRequest.GetEmail()).Return(validUser, nil).Times(1).Do(func(string) {
				defer GinkgoRecover()
			})
			authRepo.EXPECT().GetNewClaims(validUser.Email, map[string]interface{}{
				"user": validUser,
			}).Return(validClaims).Times(1).Do(func(string, map[string]interface{}) {
				defer GinkgoRecover()
			})
			authRepo.EXPECT().GetSignedToken(validClaims).Return("", errors.New(errMsg)).Times(1).Do(func(*jwt.Claims) {
					defer GinkgoRecover()
			})
			_, err := router.Login(ctx, validRequest)
			Ω(err).NotTo(BeNil())
			Ω(err.Error()).To(Equal(errMsg))
		})

		It("success", func() {
			mockToken := "1qaz2wsx"
			usersRepo.EXPECT().FindByEmail(validRequest.GetEmail()).Return(validUser, nil).Times(1).Do(func(string) {
				defer GinkgoRecover()
			})
			authRepo.EXPECT().GetNewClaims(validUser.Email, map[string]interface{}{
				"user": validUser,
			}).Return(validClaims).Times(1).Do(func(string, map[string]interface{}) {
				defer GinkgoRecover()
			})
			authRepo.EXPECT().GetSignedToken(validClaims).Return(mockToken, nil).Times(1).Do(func(claims *jwt.Claims) {
				defer GinkgoRecover()
			})
			token, err := router.Login(ctx, validRequest)
			Ω(err).To(BeNil())
			Ω(token).NotTo(BeNil())
			Ω(token.GetToken()).To(Equal(mockToken))
		})
	})
})