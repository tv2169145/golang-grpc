package repos_test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	"github.com/tv2169145/golang-grpc/types"

	. "github.com/tv2169145/golang-grpc/types"

	. "github.com/onsi/gomega"
)

var _ = Describe("UsersRepo", func() {
	var (
		err error
		usr *User

		setupData = func() {
			usr, err = NewUser(&TempUser{
				FirstName:       "Linn",
				LastName:        "Jimmy",
				Email:           "jimmy@gmail.com",
				Password:        "1234",
				ConfirmPassword: "1234",
			})
			Ω(err).To(BeNil())
		}
	)
	BeforeEach(func() {
		clearDatabase()
		setupData()
	})

	Describe("Create", func() {
		// 測試建立會員失敗
		Context("Failures", func() {
			// usr equal nil
			It("should fail with a nil user", func() {
				err = gr.Users().Create(nil)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal("validator: (nil *types.User)"))
			})
			// Invalid user data
			It("should fail with a bad user", func() {
				err = gr.Users().Create(&types.User{
					Password: usr.Password,
					Visible: true,

				})
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal("Key: 'User.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag\nKey: 'User.LastName' Error:Field validation for 'LastName' failed on the 'required' tag\nKey: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"))
			})
			It("should fail with database error", func() {
				errMsg := "database unavailable"

				mock.ExpectExec("INSERT INTO `users` (`first_name`,`last_name`,`email`,`password`,`visible`) VALUES (?, ?, ?, ?, ?)").
					WithArgs(usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Visible).
					WillReturnError(errors.New(errMsg))

				err = gr.Users().Create(usr)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
		})
		Context("Success", func() {
			It("successfully stored a user", func() {
				mock.ExpectExec("INSERT INTO `users` (`first_name`,`last_name`,`email`,`password`,`visible`) VALUES (?, ?, ?, ?, ?)").
					WithArgs(usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Visible).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err = gr.Users().Create(usr)
				Ω(err).To(BeNil())
			})
		})

	})
})
