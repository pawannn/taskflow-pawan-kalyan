package Error

type Err struct {
	Code    int
	Message string
	Data    interface{}
}

var NoErr = Err{}

func (e Err) IsEmpty() bool {
	return e.Code == 0 && e.Message == "" && e.Data == nil
}

func NewErr(code int, message string, data interface{}) Err {
	return Err{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
