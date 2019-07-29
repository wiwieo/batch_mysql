package db

import (
	"fmt"
	"strings"
	"wiwieo/batch_mysql/model"
)

func InsertDebrisType(dts []*model.DebrisType) error {
	if len(dts) == 0{
		return fmt.Errorf("there is no ad that need to persist")
	}
	insertSQL := strings.Builder{}
	insertSQL.WriteString( `insert into debris_type(title, sort) values`)
	values := "(?, ? ),"
	var params = make([]interface{}, 0, 7)
	for _, ad := range dts {
		insertSQL.WriteString(values)
		params = append(params, ad.Title)
		params = append(params, ad.Sort)
	}
	sql := insertSQL.String()
	_, err := Conn.Exec(sql[:len(sql) - 1], params...)
	if err != nil{
		panic(err)
	}
	return nil
}