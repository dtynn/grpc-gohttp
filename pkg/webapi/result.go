package webapi

// Result json result
type Result struct {
	Res struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	} `json:"res"`

	Data interface{} `json:"data"`
}
