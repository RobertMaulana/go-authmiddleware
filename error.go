package authmiddleware

const (
	ErrorNotSendRequestId = "x-request-id cannot be empty"
	ErrorNotSendUserId    = "x-token invalid / cannot be empty"
	ErrorNotSendPlatform  = "x-platform cannot be empty"

	// empty state
	emptyValue string = ""
)
