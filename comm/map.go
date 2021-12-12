package comm

import (
	"encoding/json"
	"fmt"
	"sync"
)

type DMap struct {
	sync.RWMutex
	Describe string
	dMap map[string]interface{}
	nLock int
	nRLock int
}

// 设置map初始容量
func (mp* DMap) SetCapacity(capacity int)  {
	mp.dMap = make(map[string]interface{}, capacity)
	mp.nLock = 0
	mp.nRLock = 0
}

// 设置map初始容量
func (mp* DMap) SetDescribe(des string)  {
	mp.Describe = des
}


// 取map元素个数
func (mp* DMap) GetLen() int {
	return len(mp.dMap)
}

func (mp* DMap) Get(key string) (interface{}, error) {
	mp.RLock()
	if v, ok := mp.dMap[key]; ok {
		mp.RUnlock()
		return v, nil
	} else {
		mp.RUnlock()
		return nil, fmt.Errorf("value of key[%v] not contain!", key)
	}
}

func (mp* DMap) Set(key string, value interface{}) {
	mp.Lock()
	mp.dMap[key] = value
	mp.Unlock()
}

func (mp* DMap) Del(key string) {
	mp.Lock()
	delete(mp.dMap, key)
	mp.Unlock()
}

// 设置JSON到MAP
func (mp* DMap) SetMapFromJson(data []byte) error {
	mp.Lock()
	err := json.Unmarshal(data, &mp.dMap)
	mp.Unlock()

	return err
}

//读取MAP到JSON
func (mp* DMap) GetJsonFromMap() ([]byte,error) {
	mp.RLock()
	jsons, err := json.Marshal(mp.dMap) //转换成JSON返回的是byte[]
	mp.RUnlock()
	if err != nil {
		return nil,err
	}
	return jsons, nil
}

// 遍历
func (mp* DMap) Range(f func(key string, value interface{}) bool) {
	mp.RLock()
	for k, v := range mp.dMap {
		if !f(k, v) {
			break
		}
	}
	mp.RUnlock()
}
