package msg

import "encoding/json"

type Msg struct {
	Id   string `json:"id"`
	Body string `json:"body"`
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
