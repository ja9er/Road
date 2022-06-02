package sqlmoudle

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"log"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "road"
)

//Db数据库连接池
var DB *sql.DB

type Bannerresult struct {
	Id          int64
	Task_Id     string
	Target      string
	Banner      sql.NullString
	Server      sql.NullString
	Status_Code sql.NullString
	Title       sql.NullString
	Last_time   sql.NullTime
	Pocmatch    sql.NullInt16
}

type Userinfo struct {
	Id       int64
	Username string
	Passwd   string
}

func Queryuser(username string) Userinfo {

	var user Userinfo
	sqlStr := `select * from goadmin_user where username =?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error:%v\n", err)
		return Userinfo{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Println("[-] Query Sql error:%v\n", err)
		return Userinfo{}
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&user.Id, &user.Username, &user.Passwd)
		if e != nil {
			fmt.Println(e)
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return Userinfo{}
		}

	}
	return user
}

//查询taskid获取banner匹配结果
func Queryifno() []Bannerresult {
	var banner Bannerresult
	var result []Bannerresult
	sqlStr := `select * from bigtask`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&banner.Id, &banner.Task_Id, &banner.Target, &banner.Banner, &banner.Server, &banner.Status_Code, &banner.Title, &banner.Last_time, &banner.Pocmatch)
		if e != nil {
			fmt.Println(e)
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return nil
		}
		result = append(result, banner)
	}
	return result
}

//查询taskid获取banner匹配结果
func UpdateTask(temp Bannerresult) {

	sqlStr := `UPDATE  bigtask SET Target=?,  Banner=?,  Server=?, Status_Code=?,  Title=?, Pocmatch=? where ID= ?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Prepare Sql error:\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(temp.Target, temp.Banner.String, temp.Server.String, temp.Status_Code.String, temp.Title, temp.Pocmatch.Int16, temp.Id)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Query Sql error:\n", err)
		return
	}
	defer rows.Close()

}

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8&parseTime=true"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] MySql open database fail")
		return
	}
	color.RGBStyleFromString("168,215,186").Println("[+] MySql connnect success")
}

func Insertbanner(result Bannerresult) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Insertbanner begin Tx fail")
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO task (`Task_Id`, `Target`,`Banner`,`Server`,`Status_Code`,`Title`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("[-] MySql Prepare fail")
		fmt.Println(err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Task_Id, &result.Target, &result.Banner, &result.Server, &result.Status_Code, &result.Title)
	if err != nil {
		fmt.Println(err)
		color.RGBStyleFromString("168,215,186").Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}

//查询taskid获取banner匹配结果
func QueryTask(target string) []Bannerresult {
	target = strings.Replace(target, "\r\n", "", -1)
	var banner Bannerresult
	var result []Bannerresult
	sqlStr := `select * from task where Task_Id = ?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(target)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&banner.Id, &banner.Task_Id, &banner.Target, &banner.Banner, &banner.Server, &banner.Status_Code, &banner.Title)
		if e != nil {
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return nil
		}
		result = append(result, banner)
	}
	return result
}
