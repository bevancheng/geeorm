package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

//Field represent a column of database
type Field struct {
	Name string
	Type string
	//约束条件
	Tag string
}

//Schema represent a table of database
type Schema struct {
	Model      interface{}       //被映射的对象
	Name       string            //表名
	Fields     []*Field          //字段
	FieldNames []string          //所有字段名
	fieldMap   map[string]*Field //字段名-Field映射
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTyepOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
