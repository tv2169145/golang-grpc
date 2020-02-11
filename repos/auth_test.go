package repos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tv2169145/golang-grpc/types"
)

var _ = Describe("AuthRepo", func() {
	var (
		usr *types.User
	)

	JustBeforeEach(func() {
		usr = &types.User{
			ID: 1,
			FirstName: "Linn",
			LastName: "jimmy",
			Email: "jimmy@gmail.com",
			Password: "1234",
			Visible: true,
		}
	})

	Describe("GetNewClaims", func() {
		 It("success", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": usr,
			})
			Ω(claims).NotTo(BeNil())
			Ω(claims.Set["user"]).NotTo(BeNil())
		 })
	})

	Describe("GetSignedToken", func() {
		It("should return a signed token", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": usr,
			})
			Ω(claims).NotTo(BeNil())
			Ω(claims.Set["user"]).NotTo(BeNil())
			token, err := gr.Auth().GetSignedToken(claims)
			Ω(err).To(BeNil())
			Ω(token).NotTo(BeNil())
		})
	})

	Describe("GetDataFromToken", func() {
		It("success", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": usr,
			})
			Ω(claims).NotTo(BeNil())
			Ω(claims.Set["user"]).NotTo(BeNil())
			token, err := gr.Auth().GetSignedToken(claims)
			Ω(err).To(BeNil())
			Ω(token).NotTo(BeNil())
			user, err := gr.Auth().GetDataFromToken(token)

			Ω(err).To(BeNil())
			Ω(user).NotTo(BeNil())
			Ω(user.Email).To(Equal(usr.Email))
		})

		It("should fail because token is empty", func() {
			_, err := gr.Auth().GetDataFromToken("")
			Ω(err).NotTo(BeNil())
		})

		It("should fail to return data because user is missing visible", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": map[string]interface{}{
					"id": usr.ID,
					"email": usr.Email,
				},
			})
			token, err := gr.Auth().GetSignedToken(claims)
			Ω(err).To(BeNil())
			user, err := gr.Auth().GetDataFromToken(token)
			Ω(err).NotTo(BeNil())
			Ω(user).To(BeNil())
			Ω(err.Error()).To(Equal("token is valid but user data is missing or corrupt"))
		})

		It("should fail to return data because user is missing email", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": map[string]interface{}{
					"id": usr.ID,
					"visible": true,
				},
			})
			token, err := gr.Auth().GetSignedToken(claims)
			Ω(err).To(BeNil())
			user, err := gr.Auth().GetDataFromToken(token)
			Ω(err).NotTo(BeNil())
			Ω(user).To(BeNil())
			Ω(err.Error()).To(Equal("token is valid but user data is missing or corrupt"))
		})

		It("should fail to return data because user is missing id", func() {
			claims := gr.Auth().GetNewClaims("test", map[string]interface{}{
				"user": map[string]interface{}{
					"email": usr.Email,
					"visible": true,
				},
			})
			token, err := gr.Auth().GetSignedToken(claims)
			Ω(err).To(BeNil())
			user, err := gr.Auth().GetDataFromToken(token)
			Ω(err).NotTo(BeNil())
			Ω(user).To(BeNil())
			Ω(err.Error()).To(Equal("token is valid but user data is missing or corrupt"))
		})
	})
})
