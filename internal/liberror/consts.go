package liberror

var (
	ErrorServerUnknown OurErrorCode = NewOurErrorCode("error_server_unknown")
	ErrorInterrupted   OurErrorCode = NewOurErrorCode("error_interrupted")
	ErrorNotFound      OurErrorCode = NewOurErrorCode("error_not_found")
	ErrorDataInvalid   OurErrorCode = NewOurErrorCode("error_data_invalid")
	ErrorLimitReached  OurErrorCode = NewOurErrorCode("error_limit_reached")
)
