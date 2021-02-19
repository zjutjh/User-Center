package model

type Result struct {
	Code	int			`json:"code"`
	Msg		string		`json:"message"`
	Data	interface{}	`json:"data"`
}

func NewResult(data interface{}, code int, msg ...string) *Result {
	r := &Result{Code: code, Data: data}
	if e, ok := data.(error); ok {
		if msg == nil {
			r.Msg = e.Error()
		}
	} else {
		r.Msg = "ok"
	}
	if len(msg) > 0 {
		r.Msg = msg[0]
	}
	return r
}