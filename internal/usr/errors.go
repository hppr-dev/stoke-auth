package usr

import "errors"

var (
	AuthenticationError = errors.New("Could not authenticate user")
	NoLinkedGroupsError = errors.New("No linked groups associated with user")
	UserNotFoundError   = errors.New("User not found")
	AuthSourceError     = errors.New("An error occured with the authentication source")
)
