package response

const (
	StatusOK          = 2000
	StatusParamsEmpty = 4001
	StatusFail        = 4000
)

var statusText = map[int]string{
	StatusOK:          "Success",
	StatusParamsEmpty: "参数为空",
}

// StatusText 通过状态码获取状态描述
func StatusText(code int) string {
	return statusText[code]
}
