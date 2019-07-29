package db

import (
	"fmt"
	"strings"
	"wiwieo/batch_mysql/model"
)

func InsertAd(ads []*model.Ad) error {
	if len(ads) == 0 {
		return fmt.Errorf("there is no ad that need to persist")
	}
	insertSQL := strings.Builder{}
	insertSQL.WriteString(`insert into ad(title, as_id, pic, url, addtime, sort, open, content) values`)
	values := "(?, ?, ?, ?, UNIX_TIMESTAMP(now()), ?, ?, ?),"
	var params = make([]interface{}, 0, 7)
	for _, ad := range ads {
		insertSQL.WriteString(values)
		params = append(params, ad.Title)
		params = append(params, ad.AsId)
		params = append(params, ad.Pic)
		params = append(params, ad.Url)
		params = append(params, ad.Sort)
		params = append(params, ad.Open)
		params = append(params, ad.Content)
	}
	sql := insertSQL.String()
	_, err := Conn.Exec(sql[:len(sql)-1], params...)
	if err != nil {
		panic(err)
	}
	return nil
}
