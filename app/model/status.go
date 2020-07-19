package model

const (
	InvalidFile      = "Invalid File: A .txt file is required"
	AdminRequired    = "Please contact an administrator."
	UserFound        = "User Found"
	UserNotFound     = "User Not Found"
	ClubFound        = "Club Found"
	ClubNotFound     = "Club Not Found"
	UsernameExists   = "Username already exists"
	TagCreated       = "Successfully created tag"
	TagsCreated      = "Successfully created tags"
	TagFound         = "Tag already exists"
	UsernameAlphaNum = "Username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail       = "Must be a valid email"
	EmailExists      = "Email already exists"
	FailureCode      = -1
	SuccessCode      = 1
)

type T struct{}

type Status struct {
	Code    int
	Message string
	Data    interface{}
}

type CredentialStatus struct {
	Username string
	Email    string
}

func NewStatus() *Status {
	return &Status{Code: SuccessCode, Message: "", Data: T{}}
}
