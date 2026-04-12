package requestcontext

type httpContext string

var RequestKey httpContext = "taskflowRequest"

type ReqContext struct {
	ReqID     string `json:"req_id"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}
