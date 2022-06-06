# geeorm
学习https://geektutu.com/post/geeorm.html
实现geeorm



## Day0

ORM(Object Relational Mapping)对象关系映射，通过使用描述对象和数据库之间的映射的元数据，将面向对象语言程序中的对象自动持久化到关系数据库中。

对象和数据库的映射：
表(table)->类(class/struct)
记录(record，row)->对象(object)
字段(field,column)->对象属性(attribute)

```sql
CREATE TABLE `User` (`Name` text, `Age` integer);
INSERT INTO `User` (`Name`, `Age`) VALUES ("Tom", 18);
SELECT * FROM `User`;
```
用ORM框架为
```go
type User struct{
    Name string
    Age int
}
orm.CreateTable(&User{})
orm.Save(&User{"Tom",18})
var users []User
orm.Find(&users)
```


```go
type Account struct{
    Username string
    Password string
}
orm.CreateTable(&Account{})
```

那么如何根据任意类型的指针，得到其对应的结构体的信息？
这就要用到GO的反射机制(reflect)，通过反射，可以得到对象对应的结构体名称，成员变量，风法等信息。
```go
typ := reflect.Indirect(reflect.ValueOf(&Account{})).Type()
fmt.Println(typ.Name())

for i:=0 ; i<typ.NumField();i++{
    field := typ.Field(i)
    fmt.Println(field.name)
}
```
`reflect.ValueOf`获取指针对应的反射值
`reflect.Indirect`获取指针指向的对象的反射值
`(reflect.Type).Name()`返回类名(字符串)
`(reflect.Type).Field(i)`获取第i个成员变量

设计一个ORM框架，需要考虑功能特性的优先级。比较广泛使用的ORM框架有GORM和XORM。除了基础功能，如表的操作、记录的增删改查，gorm还实现了关联关系、回调插件；xorm实现了读写分离、数据同步、导入导出。
geeorm设计参考xorm，细节实现参考gorm。

## Day1 database/sql

sql.Open()连接数据库，_import导入时会注册sqlite3的驱动。
Exec()执行SQL语句
Exec(),Query(),QueryRow()接收1或多个入参，第一个入参时SQL语句，后边是占位符？对应的值，占位符防止SQL注入

首先实现一个log库，这个简易的 log 库具备以下特性：

- 支持日志分级（Info、Error、Disabled 三级）。
- 不同层级日志显示时使用不同的颜色区分。
- 显示打印日志代码对应的文件名和行号。

log.Lshortfile 支持显示文件名和代码行号。

NewEngine中首先与数据库建立连接(Open),之后Ping看是否连接接通

## Day2 对象表结构映射
  
``` go
type Dialect interface {
	DataTyepOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}
```

`DataTyepOf`将GO语言中的类型转换为该数据的数据类型
`TableExistSQL`某个表是否存在SQL语句
`RegisterDialect`用于注册Dialect实例，如果新增对某个数据库的支持，那么调用注册可注册到全局

在数据库中创建一个表需要的要素：
- 表名(table name)----结构体名(struct name)
- 字段名和字段类型----成员变量和类型
- 额外约束(非空、主键)----成员变量的Tag

```go
type User struct{
    Name string `geeorm:"PRIMARY KEY"`
    Age int
}
```

```sql
CREATE TABLE `User` (`Name` text PRIMARY KEY, `Age` integer)
```

## Day3
## Day4

1，支持Update、Delete、Count
record.go中Update方法支持两种入参，平铺的kv对和map类型的键值对。generator接受的参数必须是map类型键值对，因此Update需要自动转换为map类型。
2，链式调用
``` go
s := geeorm.NewEngine("sqlite3", "gee.db").NewSession()
var users []User
s.Where("Age > 18").Limit(3).Find(&users)
```
## Day5
Hook，提前在可能增加功能的地方埋好一个钩子，当需要重新修改或增加这个地方的逻辑，把扩展类或方法挂载到这个点即可。
如GitHub支持的travis持续集成服务，当git push时，会触发travis拉取新的代码进行构建。
如IDE中的 Ctrl+S后自动格式化代码。前端hot reload机制。
对于ORM来说，扩展点在CURD前后比较合适。

