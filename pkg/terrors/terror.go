package terrors

type Terror struct {
	ErrCode        string `json:"code"`
	HttpStatusCode int    `json:"-"`
	ErrMessage     string `json:"message"`
}

func (e *Terror) Error() string {

	return e.ErrMessage
}
