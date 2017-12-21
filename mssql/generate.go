package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"fmt"
	"os"
	"strings"
)

type Table struct{
	ID int `db:"object_id"`
	Name string  `db:"name"`
}
type Column struct{
	ID int `db:"object_id"`
	Name string  `db:"name"`
	Type string `db:"type"`
	IsNull bool `db:"is_nullable"`
	IsIdentity bool  `db:"is_identity"`
	Description string `db:"value"`
}

func Generator(connection string) (err error) {
	db, err := sqlx.Connect("mssql", connection)
	if err != nil {
		return
	}
	defer db.Close()
	var tables []Table
	err = db.Select(&tables, "select object_id,name from sys.tables where type='U'")
	if err != nil {
		return
	}
	for _, table := range tables {
		var columns []Column
		var sql=fmt.Sprintf(`select object_id,c.name as name,t.name as type,c.is_nullable as is_nullable,c.is_identity as is_identity,isnull(des.value,'') as value from sys.columns c
		left join sys.extended_properties  des on c.column_id=des.minor_id and c.object_id=des.major_id
		left join sys.types t on c.system_type_id=t.system_type_id
		where object_id = %d`,table.ID)
		err = db.Select(&columns, sql)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(columns)>0 {
			WriteToFile(table.Name,columns)
		}
	}
	return nil
}

func WriteToFile(tableName string,columns []Column) {
	dir := "files/"
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return
		}
	}
	fileName := dir + tableName + ".go"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open file error !", err)
		return
	}
	defer file.Close()
	template := fmt.Sprintf("packet model\r\n\r\ntype %v struct{\r\n", tableName)
	for _, column := range columns {
		template = template + fmt.Sprintf("\t%v\t%v\t`db:\"%v\"`\t%v\r\n", CheckColumnName(column.Name), CheckSQLType(column.Type), column.Name, CheckDes(column.Description))
	}
	template = template + "}"
	file.WriteString(template)
}
//首字母大写
func CheckColumnName(name string) string {
	if len(name)==0{
		return ""
	}
	return strings.ToUpper(name[0:1]) + name[1:]
}
//检查类型
func CheckSQLType(tp string) string {
	switch tp {
	case "int":
		return "int32"
	case "bit":
		return "bool"
	case "bigint":
		return "int64"
	case "datetime":
		return "time.Time"
	case "date":
		return "time.Time"
	case "datetime2":
		return "time.Time"
	case "decimal":
		return "float64"
	case "float":
		return "float64"
	case "varchar":
		return "string"
	case "nvarchar":
		return "string"
	case "text":
		return "string"
	case "tinyint":
		return "uint"
	default:
		return "string"
	}
}

func CheckDes(des string)string{
	if des==""{
		return ""
	}
	return `//`+des
}





