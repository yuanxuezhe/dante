package msg

import (
	"dante/core/log"
	"encoding/json"
)

type Msg struct {
	Id   string `json:"id"`
	Mac  string `json:"mac"`
	Body string `json:"body"`
}

type Result struct {
	Module string `json:"module"` // 模块类型
	Status string `json:"status"` // 状态
	Code   int    `json:"code"`   // 错误码
	Msg    string `json:"msg"`    // 消息
	Data   string `json:"data"`   // 结果
}

func PackageMsg(id string, body string) []byte {
	m := &Msg{
		Id:   id,
		Body: body,
	}

	jsons, err := json.Marshal(m) //转换成JSON返回的是byte[]

	if err != nil {
		panic(err)
	}
	return jsons
}

//
func ResultPackege(moduleType string, code int, msg string, data interface{}) []byte {
	result := &Result{}
	if code == 0 {
		result.Status = "ok"
	} else {
		result.Status = "err"
	}

	result.Module = moduleType

	result.Code = code

	result.Msg = msg

	buff, _ := json.Marshal(data)
	result.Data = string(buff)

	resbuff, _ := json.Marshal(result)

	if result.Status == "ok" {
		log.Release("[%10s]%s", moduleType, string(resbuff))
	} else {
		log.Error("[%10s]%s", moduleType, string(resbuff))
	}

	return resbuff
}
