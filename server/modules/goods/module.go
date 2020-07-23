package goods

import (
	"encoding/json"
	"gitee.com/yuanxuezhe/dante/core/module"
	"gitee.com/yuanxuezhe/dante/core/module/base"
	. "gitee.com/yuanxuezhe/dante/core/msg"
	"gitee.com/yuanxuezhe/dante/server/tables"
	"sync"
)

type GoodsInfo struct {
	Goodsid   int32  `json:"goodsid"`   //编号
	Goodsname string `json:"goodsname"` //名称
	Type      int    `json:"type"`      //商品类型
	Source    string `json:"source"`    //来源
	Url       string `json:"url"`       //链接
	Imgurl    string `json:"imgurl"`    //图片链接
	Brand     int    `json:"brand"`     //品牌
	Status    int    `json:"status"`    //状态
	Date      int    `json:"date"`      //日期
	Time      int    `json:"time"`      //时间
}

var NewModule = func() module.Module {
	mod := &GoodsManage{Basemodule: base.Basemodule{ModuleType: "Goods", ModuleVersion: "1.2.9"}}
	mod.Basemodule.DoWork = mod.DoWork
	return mod
}

type GoodsManage struct {
	base.Basemodule
	rw sync.RWMutex
}

func (m *GoodsManage) DoWork(buff []byte) ([]byte, error) {
	var err error
	// 解析收到的消息
	t_goods := tables.Goods{}
	json.Unmarshal(buff, t_goods)
	if err != nil {
		return nil, err
	}

	data, err := t_goods.QueryByStatus()
	if err != nil {
		return nil, err
	}
	return ResultPackege(m.ModuleType, 0, "Get goodsinfo successful!", data), nil
}
