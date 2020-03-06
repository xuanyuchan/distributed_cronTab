package Api

import "encoding/json"

type apiResponse struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func buildResponse(errno int, errMsg string, data interface{}) ([]byte, error) {
	resp := &apiResponse{
		ErrNo:  errno,
		ErrMsg: errMsg,
		Data:   data,
	}
	return json.Marshal(resp)
}
