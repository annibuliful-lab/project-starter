package error_utils

import (
	"errors"
)

var (
	InternalServerError     = errors.New("internal server error")
	SignTokenFailed         = errors.New("sign token is failed")
	TokenIsInvalid          = errors.New("Token is invalid")
	ClaimsIsInvalid         = errors.New("Claims is invalid")
	TokenExpire             = errors.New("token has expire")
	TokenVerificationFailed = errors.New("token verification failed")
	ErrInvalidHash          = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion  = errors.New("incompatible version of argon2")
	ForbiddenOperation      = errors.New("forbidden operation")
	ContactOwner            = errors.New("please contact project owner")
	DataConflict            = errors.New("data conflict")
	IdNotfound              = errors.New("id not found")
	Notfound                = errors.New("data not found")
	ErrInvalidIMSI          = errors.New("invalid Imsi")
	ErrInvalidIMEI          = errors.New("invalid Imei")
)
