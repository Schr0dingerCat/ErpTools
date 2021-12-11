package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type JsHttp struct {
	Url  string `json:"url"`
	Port int    `json:"port"`
}

type JsMssql struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type JsConfig struct {
	Http  JsHttp  `json:"http"`
	Mssql JsMssql `json:"mssql"`
}

type JsPrdno struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

var (
	Sqlconn *sql.DB
	Config  JsConfig
	Options []JsPrdno
	err error
)

func main() {
	err = LoadJsonFile("./json/config.json", &Config)
	if err != nil {
		fmt.Print("加载配置文件config.json错误：")
		fmt.Println(err)
		return
	} else {
		if Config.Http.Port <= 0 || Config.Http.Port > 65535 {
			Config.Http.Port = 3000
		}
	}

	InitDB()
	defer Sqlconn.Close()

	r := gin.Default()
	r.Static("/assets", "./dist/assets")
	r.StaticFile("erptools", "./dist/index.html")
	// r.GET("/about.html", GetAbout)
	// r.GET("/erptools", GetErpTools)
	r.POST("/erptools", PostErpTools)

	r.Run(fmt.Sprintf("%s:%d", Config.Http.Url, Config.Http.Port))
}

// 初始化数据库连接
func InitDB() {
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", Config.Mssql.Server, Config.Mssql.Port, Config.Mssql.Database, Config.Mssql.User, Config.Mssql.Password)
	Sqlconn, err = sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("open database connection failed: ", err.Error())
	}
	Sqlconn.SetMaxOpenConns(2000)
	Sqlconn.SetMaxIdleConns(1000)
}

// 显示about.html页
func GetAbout(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.HTML(http.StatusOK, "./dist/about.html", gin.H{
		"config": Config,
	})
}

// 测试使用
func GetErpTools(c *gin.Context) {
	cmd := c.DefaultQuery("cmd", "null")
	switch cmd {
	case "getprdno":
		c.JSON(http.StatusOK, gin.H{
			"cmd":     cmd,
			"message": "获取货品代号",
			"error":   0,
		})
	case "config.json":
		c.JSON(http.StatusOK, Config)
	default:
		c.HTML(http.StatusOK, "./dist/about.html", gin.H{
			"config": Config,
		})
	}
}

// 数据获取处理接口
func PostErpTools(c *gin.Context) {
	// cmd := c.DefaultQuery("cmd", "null")
	jsonMap := make(map[string]string)
	c.BindJSON(&jsonMap)
	switch jsonMap["cmd"] {
	case "getprdno": //加载页面 /erptools时，获取后台数据
		// 获取prdno列表数据
		Options = GetPrdnoFromDB()
		// 返回数据
		c.JSON(http.StatusOK, gin.H{
			"options": Options,
			"message": "获取货品代号",
			"error":   0,
		})
	default:
		c.JSON(http.StatusFound, gin.H{
			"message": "Error 302",
			"error":   302,
		})
	}
}

func GetPrdnoFromDB() []JsPrdno {
	s := `SELECT PRD_NO FROM PRDT WHERE (KND = '2') AND (IDX1 = '11' OR IDX1 = '12' OR IDX1 = '13' OR IDX1 = '14' OR IDX1 = '15' OR IDX1 = '16' OR IDX1 = '17' OR IDX1 = '18' OR IDX1 = '19') ORDER BY PRD_NO`
	rows, err := Sqlconn.Query(s)
	if err != nil {
		log.Fatal("GetPrdnoFromDB rows err: ", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var prdno string
		rows.Scan(&prdno)
		Options = append(Options, JsPrdno{prdno, prdno})
	}
	// Options = []JsPrdno{{"Q-001"}, {"Q-002"}, {"Q-003"}, {"Q-004"}, {"Q-005"}, {"Q-006"}, {"Q-007"}, {"Q-008"}, {"Q-009"}, {"Q-010"}, {"Q-011"}, {"Q-012"}, {"Q-013"}, {"Q-014"}, {"Q-015"}, {"Q-016"}, {"Q-017"}, {"Q-018"}}
	return Options
}

func LoadFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, err
}

func LoadJsonString(content []byte, obj interface{}) error {
	err := json.Unmarshal(content, obj)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func LoadJsonFile(path string, obj interface{}) error {
	content, err := LoadFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return LoadJsonString(content, obj)
}
