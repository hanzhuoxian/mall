package userv1

import "github.com/go-playground/validator/v10"

var v = validator.New()

// Validate checks that the request contains a user identifier.
func (r *GetUserRequest) Validate() error {
	return v.Var(r.GetInstanceId(), "required")
}

// Validate checks the required fields of the create request.
func (r *CreateUserRequest) Validate() error {
	type req struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}
	return v.Struct(req{
		Username: r.GetUsername(),
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
	})
}

// Validate checks that the request contains a user identifier.
func (r *UpdateUserRequest) Validate() error {
	return v.Var(r.GetInstanceId(), "required")
}

// Validate checks that the request contains a user identifier.
func (r *DeleteUserRequest) Validate() error {
	return v.Var(r.GetInstanceId(), "required")
}

// Validate checks that the request contains at least one user identifier.
func (r *DeleteCollectionRequest) Validate() error {
	return v.Var(r.GetInstanceIds(), "required,min=1")
}

// Validate checks the required credentials in the request.
func (r *AuthenticateUserRequest) Validate() error {
	type req struct {
		Identifier string `validate:"required"`
		Password   string `validate:"required"`
	}
	return v.Struct(req{
		Identifier: r.GetIdentifier(),
		Password:   r.GetPassword(),
	})
}
