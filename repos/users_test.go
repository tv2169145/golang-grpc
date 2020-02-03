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
				err := gr.Users().Create(nil)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal("validator: (nil *types.User)"))
			})
			// Invalid user data
			It("should fail with a bad user", func() {
				err := gr.Users().Create(&types.User{
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

				err := gr.Users().Create(usr)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
		})
		Context("Success", func() {
			It("successfully stored a user", func() {
				mock.ExpectExec("INSERT INTO `users` (`first_name`,`last_name`,`email`,`password`,`visible`) VALUES (?, ?, ?, ?, ?)").
					WithArgs(usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Visible).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := gr.Users().Create(usr)
				Ω(err).To(BeNil())
			})
		})
	})

	Describe("FindById", func() {
		Context("Failures", func() {
			It("should fail with a bad id", func() {
				_, err := gr.Users().FindById(0)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal("invalid user id"))
			})
			It("should fail with database error", func() {
				usr.ID = 1
				errMsg := "database unavailable"
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `id`=? LIMIT 1").
					WithArgs(usr.ID).
					WillReturnError(errors.New(errMsg))
				_, err := gr.Users().FindById(usr.ID)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
			It("should fail with not found a user", func() {
				errMsg := "unable to find a user"
				usr.ID = 1
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `id`=? LIMIT 1").
					WithArgs(usr.ID).
					WillReturnRows(sqlmock.NewRows([]string{}))
				_, err := gr.Users().FindById(usr.ID)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
		})

		Context("success", func() {
			It("should successfully find a user", func() {
				usr.ID = 1
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `id`=? LIMIT 1").
					WithArgs(usr.ID).
					WillReturnRows(sqlmock.NewRows(
						[]string{"id", "first_name", "last_name", "email", "password", "visible"},
					).AddRow(usr.ID, usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Visible))
				foundUser, err := gr.Users().FindById(usr.ID)
				Ω(err).To(BeNil())
				Ω(foundUser).To(BeEquivalentTo(usr))
			})
		})
	})

	Describe("FindByEmail", func() {
		Context("failures", func() {
			It("should fail with empty email", func() {
				_, err := gr.Users().FindByEmail("")
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal("invalid user email"))
			})
			It("should fail with database error", func() {
				errMsg := "database unavailable"
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `email`=? LIMIT 1").
					WithArgs(usr.Email).
					WillReturnError(errors.New(errMsg))
				_, err := gr.Users().FindByEmail(usr.Email)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
			It("should fail with not found a user", func() {
				errMsg := "unable to find a user"
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `email`=? LIMIT 1").
					WithArgs(usr.Email).
					WillReturnRows(sqlmock.NewRows([]string{}))
				_, err := gr.Users().FindByEmail(usr.Email)
				Ω(err).NotTo(BeNil())
				Ω(err.Error()).To(Equal(errMsg))
			})
		})

		Context("success", func() {
			It("should successfully find a user", func() {
				mock.ExpectQuery("SELECT `id`, `first_name`, `last_name`, `email`, `password`, `visible` FROM `users` WHERE `email`=? LIMIT 1").
					WithArgs(usr.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "visible"}).AddRow(usr.ID, usr.FirstName, usr.LastName, usr.Email, usr.Password, usr.Visible))
				foundUser, err := gr.Users().FindByEmail(usr.Email)
				Ω(err).To(BeNil())
				Ω(foundUser).To(BeEquivalentTo(usr))
			})
		})
	})
})
