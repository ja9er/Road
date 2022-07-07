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
	UserName = "root"
	Password = "root"
	Ip       = "10.10.10.166"
	Port     = "3306"
	DbName   = "road"
)

//const (
//	UserName = "root"
//	Password = "Docker@mysql123"
//	Ip       = "139.155.75.156"
//	Port     = "3306"
//	DbName   = "banner"
//)

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
	Status   string
}

type Linkinfo struct {
	Uuid        string
	Ip_addr     string
	Update_time sql.NullTime
	Banner      string
	Success     int
	Fail        int
}
type Linktask struct {
	Task_Id       string
	UUID          string
	Update_time   sql.NullTime
	Update_Order  string
	Update_Result string
	Update_Status int
}

type Taskjob struct {
	Id        int64
	Task_Id   string
	Fofaquery string
	Progress  int64
}

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{UserName, ":", Password, "@tcp(", Ip, ":", Port, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
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

func Queryuser(username string) []Userinfo {
	var users []Userinfo
	var user Userinfo
	sqlStr := `select * from goadmin_user where username like ?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		log.Println("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&user.Id, &user.Username, &user.Passwd, &user.Status)
		if e != nil {
			fmt.Println(e)
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return nil
		}
		users = append(users, user)
	}
	return users
}

func QueryPOCmatch(target int) []Bannerresult {
	var banner Bannerresult
	var result []Bannerresult
	sqlStr := `select * from bigtask where Pocmatch = ?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(target)
	if err != nil {
		log.Println("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&banner.Id, &banner.Task_Id, &banner.Target, &banner.Banner, &banner.Server, &banner.Status_Code, &banner.Title, &banner.Last_time, &banner.Pocmatch)
		if e != nil {
			log.Println("[-] read DataBase error: ", e)
			return nil
		}
		result = append(result, banner)
	}
	return result
}
func Queryinfo(sqlStr string) []Bannerresult {
	var banner Bannerresult
	var result []Bannerresult
	//sqlStr := `SELECT * FROM bigtask WHERE   Title = "UniFi Video"`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		color.RGBStyleFromString("168,215,186").Println("[-] Query Sql error:\n", err)
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
func Queryifno() []Bannerresult {
	var banner Bannerresult
	var result []Bannerresult
	sqlStr := `select * from  bigtask`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println("[-] Query Sql error:%v\n", err)
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

func Queryconnectinfo() []Linkinfo {
	var link Linkinfo
	var result []Linkinfo
	sqlStr := `select * from  connect_equipment`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&link.Uuid, &link.Ip_addr, &link.Update_time, &link.Banner, &link.Success, &link.Fail)
		if e != nil {
			fmt.Println(e)
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return nil
		}
		result = append(result, link)
	}
	return result
}
func DeleteconnectEquipment(result Linkinfo) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("delete from connect_equipment where UUID =?")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Uuid)
	if err != nil {
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}

func Queryconnecttaskinfo(uuid string) []Linktask {
	var link Linktask
	var result []Linktask
	sqlStr := `select * from  equipment_task_info where uuid=? order by UPDATE_TIME  DESC`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Printf("[-] Prepare Sql error:%v\n", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(uuid)
	if err != nil {
		log.Printf("[-] Query Sql error:%v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&link.Task_Id, &link.UUID, &link.Update_time, &link.Update_Order, &link.Update_Result, &link.Update_Status)
		if e != nil {
			fmt.Println(e)
			color.RGBStyleFromString("168,215,186").Println("[-] read DataBase error")
			return nil
		}
		result = append(result, link)
	}
	return result
}
func Insertconnecttask(result Linktask) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO equipment_task_info (`Task_Id`, `UUID`,`UPDATE_ORDER`,`UPDATE_RESULT`,`UPDATE_status`) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Task_Id, &result.UUID, &result.Update_Order, &result.Update_Result, &result.Update_Status)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}
func Deleteconnecttask(result Linktask) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("delete from equipment_task_info where Task_Id =?")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Task_Id)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}

func Queryconnnectcount() []int64 {
	var connectnumber []int64
	var tempnumber int64
	sqlStr := `SELECT COUNT(*)  FROM connect_equipment WHERE  datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 6 DAY ))=0 union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 5 DAY ))=0 union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 4 DAY ))=0 union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 3 DAY ))=0 union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 2 DAY ))=0 union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,date_sub( curdate( ), INTERVAL 1  DAY ))=0  union all SELECT COUNT(*)  FROM connect_equipment WHERE datediff(UPDATE_TIME,CURRENT_DATE())=0`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error: ", err)
		return connectnumber
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println("[-] Query Sql error: ", err)
		return connectnumber
	}
	defer rows.Close()

	for rows.Next() {
		e := rows.Scan(&tempnumber)
		if e != nil {
			log.Println("[-] read DataBase error: ", e)
			return connectnumber
		}
		connectnumber = append(connectnumber, tempnumber)
	}
	return connectnumber
}

//scan task
func Querytaskinfo() []Taskjob {
	var task Taskjob
	var result []Taskjob
	sqlStr := `select * from  taskmanager`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare Sql error: ", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println("[-] Query Sql error: ", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		e := rows.Scan(&task.Id, &task.Task_Id, &task.Progress, &task.Fofaquery)
		if e != nil {
			log.Println("[-] read DataBase error: ", e)
			return nil
		}
		result = append(result, task)
	}
	return result
}
func Inserttask(result Taskjob) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO taskmanager (`Task_Id`,`Fofaquery`, `Progress`) VALUES (?,?,?)")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Task_Id, &result.Fofaquery, &result.Progress)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}
func Upadtetask(result Taskjob) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("UPDATE taskmanager set  `Progress` =? where Task_id=?")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Progress, &result.Task_Id)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}
func Deletetask(result Taskjob) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("delete from taskmanager where Id =?")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Id)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}

func Deletebannerfromtask(Taskid string) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("delete from task where task_id =?")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(Taskid)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}
func Insertbanner(result Bannerresult) bool {
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Println("[-] Insertbanner begin Tx fail", err)
		return false
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO task (`Task_Id`, `Target`,`Banner`,`Server`,`Status_Code`,`Title`,`POcmatch`) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("[-] MySql Prepare fail: ", err)
		return false
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(&result.Task_Id, &result.Target, &result.Banner, &result.Server, &result.Status_Code, &result.Title, &result.Pocmatch)
	if err != nil {
		fmt.Println(err)
		log.Println("[-] MySql Exec fail", err)
		return false
	}
	//将事务提交
	tx.Commit()
	return true
}

func DeleteTask(ID int64) {
	sqlStr := `DELETE FROM  task where ID= ?`
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		log.Println("[-] Prepare DELETE Sql error: ", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(ID)
	if err != nil {
		log.Println("[-] Query DELETE Sql error:\n", err)
		return
	}
	defer rows.Close()
}

//查询taskid获取banner匹配结果
func UpdateTask(temp Bannerresult) {
	sqlStr := `UPDATE  task SET Target=?,  Banner=?,  Server=?, Status_Code=?,  Title=?, Pocmatch=? where ID= ?`
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
