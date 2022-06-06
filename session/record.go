package session

import (
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value)) //

	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	results, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return results.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()                                   //获取切片单个元素的类型destType
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable() //New()创建新的destType实例，做Model入参，映射表结构

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames) //clause构造出SELECT语句
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows() //查询所有符合条件的记录rows
	if err != nil {
		return err
	}

	for rows.Next() { //遍历每一行
		dest := reflect.New(destType).Elem() //创建实例，字段平铺
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface()) //将字段每一列值赋值给values切片
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest)) //dest添加到切片destSlice中
	}
	return rows.Close()
}
