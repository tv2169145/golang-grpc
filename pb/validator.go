package grpc

import (
	vCop "gopkg.in/go-playground/validator.v9"
)

var validator *vCop.Validate

func init() {
	validator = vCop.New()

	// Users
	validator.RegisterStructValidation(func(sl vCop.StructLevel) {
		r := sl.Current().Interface().(CreateUserRequest)

		if r.GetNewUser() == nil {
			sl.ReportError("NewUser", "NewUser", "NewUser", "valid-newUser", "")
		} else {
			if len(r.GetNewUser().GetEmail()) == 0 {
				sl.ReportError("Email", "email", "Email", "valid-email", "")
			}
			if len(r.GetNewUser().GetFirstName()) == 0 {
				sl.ReportError("FirstName", "firstName", "FirstName", "valid-firstName", "")
			}
			if len(r.GetNewUser().GetFirstName()) == 0 {
				sl.ReportError("LastName", "lastName", "LastName", "valid-lastName", "")
			}
			if len(r.GetNewUser().GetPassword()) == 0 {
				sl.ReportError("Password", "password", "Password", "valid-password", "")
			}
			if len(r.GetNewUser().GetConfirmPassword()) == 0 {
				sl.ReportError("ConfirmPassword", "confirmPassword", "ConfirmPassword", "valid-confirmPassword", "")
			}
		}

	}, CreateUserRequest{})
}

func Validate(t interface{}) error {
	return validator.Struct(t)
}






