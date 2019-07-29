package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"wiwieo/batch_mysql/db"
	"wiwieo/batch_mysql/model"
)

type DebrisType int

func (DebrisType) WriteToMySQL(contents [][]byte) error {
	var dts = make([]*model.DebrisType, 0, len(contents))
	for _, b := range contents {
		var dt *model.DebrisType
		err := json.Unmarshal(b, &dt)
		if err != nil {
			return err
		}
		dts = append(dts, dt)
	}
	glog.V(3).Infoln(fmt.Sprintf("此次共有%d条数据插入表【debris_type】", len(dts)))
	return db.InsertDebrisType(dts)
}
