package myerrors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
	ErrCodeTokenCreation   ErrorCode = "TOKEN_CREATION"
	ErrCodeDatabase        ErrorCode = "DATABASE"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenInvalid  = errors.New("refresh token invalid or expired")
)

const (
	AccessTokenCreateErrorMessage         = "Error creating access token."
	RefreshTokenCreateErrorMessage        = "Error creating refresh token."
	RefreshTokenCreateInDBErrorMessage    = "Error creating refresh token in DB."
	RefreshTokenCreateInRedisErrorMessage = "Error creating refresh token in redis."
	AuthorizationHeaderEmptyErrorMessage  = "Authorization header is empty."
	InvalidTokenErrorMessage              = "Token is invalid."
	InvalidBearerTokenFormatErrorMessage  = "Invalid bearer token format."
	InvalidSignatureAlgorithmErrorMessage = "Invalid signature algorithm."
	ParseTokenErrorMessage                = "Error parsing token."
	CheckUserExistsByEmailErrorMessage    = "Error checking user exists by email."
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code ErrorCode, message string, err error) AppError {
	return AppError{Code: code, Message: message, Err: err}
}

// Хелперы для создания конкретных ошибок
func NewUnauthorizedErr(message string, err error) AppError {
	return NewAppError(ErrCodeUnauthorized, message, err)
}

func NewTooManyRequestsErr(message string, err error) AppError {
	return NewAppError(ErrCodeTooManyRequests, message, err)
}

func NewTokenErr(message string, err error) AppError {
	return NewAppError(ErrCodeTokenCreation, message, err)
}

func NewDatabaseErr(message string, err error) AppError {
	return NewAppError(ErrCodeDatabase, message, err)
}

func NewRepositoryErr(message string, err error) AppError {
	return NewAppError(ErrCodeDatabase, message, err)
}
