package model

const (
	FailureManagerRemove   = "Unable to remove manager"
	FailureManagerAddition = "Unable to add manager"
	SuccessManagerRemove   = "Successfully removed manager"
	SuccessManagerAddition = "Successfully added manager"

	TagUpdateError = "Error Updating Tag"
	TagUpdated     = "Tag Updated"
	TagCreated     = "Successfully created tag"
	TagsCreated    = "Successfully created tags"
	TagExists      = "Tag already exists"
	TagsUpdated    = "Tags Updated"
	TagsFound      = "Tags Found"
	TagNotFound    = "Tag doesn't exist"

	UserFound     = "User Found"
	UserNotFound  = "User Not Found"
	ClubsFound    = "Clubs Found"
	ClubsNotFound = "Clubs Not Found"
	ClubFound     = "Club Found"
	ClubNotFound  = "Club Not Found"

	UsernameExists   = "Username already exists"
	UsernameAlphaNum = "Username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail       = "Must be a valid email"
	EmailExists      = "Email already exists"

	AdminRequired = "Please contact an administrator."
	InvalidFile   = "Invalid File: A .txt file is required"
	FailureCode   = -1
	SuccessCode   = 1

	NotApproved  = "Sorry, your account has not been approved yet"
	LoginSuccess = "Successfully logged in"
	LoginFailure = "Username or Password is Incorrect"
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
