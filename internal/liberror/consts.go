package liberror

var (
	ErrorServerUnknown OurErrorCode = NewOurErrorCode("error_server_unknown")
	ErrorInterrupted   OurErrorCode = NewOurErrorCode("error_interrupted")
	ErrorNotFound      OurErrorCode = NewOurErrorCode("error_not_found")
	ErrorDataInvalid   OurErrorCode = NewOurErrorCode("error_data_invalid")
	ErrorLimitReached  OurErrorCode = NewOurErrorCode("error_limit_reached")
)

var (
	ErrorHttpStatusCode500 OurErrorCode = NewOurErrorCode("error_http_code_500")
	ErrorHttpStatusCode300 OurErrorCode = NewOurErrorCode("error_http_code_300")
	ErrorHttpStatusCode400 OurErrorCode = NewOurErrorCode("error_http_code_400")
)
