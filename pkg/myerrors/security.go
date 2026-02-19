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
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrParseData           ErrorCode = "PARSE_ERROR"
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
	UserDoesNotExistErrorMessage          = "Такого пользователя не сушествует. Убедитесь что ввели номер телефона правильно."
	ParsingDateErrorMessage               = "Ошибка обработки даты. Убедитесь что формат даты соответсвует ДД:ММ:ГГГГ."
)

// Response — стандартный ответ с ошибкой
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
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

// ToResponse преобразует AppError в Response для отправки клиенту
func (e AppError) ToResponse() Response {
	r := Response{
		Message: e.Message,
	}
	if e.Err != nil {
		r.Error = e.Err.Error()
	}
	return r
}

func NewAppError(code ErrorCode, message string, err error) AppError {
	return AppError{Code: code, Message: message, Err: err}
}

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

func NewValidationError(message string, err error) AppError {
	return NewAppError(ErrCodeValidation, message, err)
}

func NewParseErr(message string, err error) AppError {
	return NewAppError(ErrParseData, message, err)
}
