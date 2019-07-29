package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"wiwieo/batch_mysql/db"
	"wiwieo/batch_mysql/model"
)

type Ad int

func (Ad) WriteToMySQL(contents [][]byte) error {
	var ads = make([]*model.Ad, 0, len(contents))
	err := json.Unmarshal(TwoToOne(contents), &ads)
	if err != nil {
		return err
	}
	glog.V(3).Infoln(fmt.Sprintf("即将有%d条数据插入表【ad】", len(ads)))
	return db.InsertAd(ads)
}

//func (Ad) WriteToMySQL(contents [][]byte) error {
//	var ads = make([]*model.Ad, 0, len(contents))
//	for _, b := range contents {
//		var ad *model.Ad
//		err := json.Unmarshal(b, &ad)
//		if err != nil {
//			return err
//		}
//		ads = append(ads, ad)
//	}
//	glog.V(3).Infoln(fmt.Sprintf("即将有%d条数据插入表【ad】", len(ads)))
//	return db.InsertAd(ads)
//}
