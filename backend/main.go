package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"erptools/opencc"
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

type JsFormSo struct {
	Prdno string `json:"prdno"`
	Dateso []string `json:"dateso"`
}

type JsSoData struct {
	Sono string `json:"sono"`
	EstItmSo int `json:"estitmso"`
	Prdno string `json:"prdno"`
	Prdname string `json:"prdname"`
	QtySo float32 `json:"qtyso"`
	Cusname string `json:"cusname"`
	Estdd string `json:"estdd"`
	ClsMpId string `json:"clsmpid"`
	Mono string `json:"mono"`
	QtySoLj float32 `json:"qtysolj"`
	BilType string `json:"biltype"`
	Status string `json:"status"`
}

type JsTzData struct {
	Tzno string `json:"tzno"`
	Depname string `json:"depname"`
	Zcname string `json:"zcname"`
	Qty float32 `json:"qty"`
	Qtyfin float32 `json:"qtyfin"`
	Qtylost float32 `json:"qtylost"`
	Qtybf float32 `json:"qtybf"`
	Qtysy float32 `json:"qtysy"`
	Qtypgs float32 `json:"qtypgs"`
	Mydinge float32 `json:"mydinge"`
}

type JsBcpSData struct {
	Prdno string `json:"prdno"`
	Zcname string `json:"zcname"`
	Qty float32 `json:"qty"`
}

var (
	tw2s, _ = opencc.New("tw2s")
	Sqlconn *sql.DB
	Config  JsConfig
	Options []JsPrdno
	err error
)

func main() {

	err = LoadJsonFile("./config/config.json", &Config)
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
	case "getsolist":
		// 处理函数
		SoDatas := GetSoTableFromDB(jsonMap["args"])
		// 返回数据
		c.JSON(http.StatusOK, gin.H{
			"sodatas": SoDatas,
			"message": "获取受订数据",
			"error":   0,
		})
	case "gettzlist":
		// 获取通知明细
		TzDatas, BcpsDatas, Qtycp := GetTzTableFromDB(jsonMap["mono"], jsonMap["prdno"])
		c.JSON(http.StatusOK, gin.H{
			"tzdatas": TzDatas,
			"qtycp": Qtycp,
			"bcpsdatas": BcpsDatas,
			"message": "获取通知单数据",
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

func GetSoTableFromDB(JsFormSoString string) []JsSoData {
	var sodatas []JsSoData
	// 获取参数
	formSO := JsFormSo{}
	_ = json.Unmarshal([]byte(JsFormSoString), &formSO)
	// utc时间转本地时间 
	t1, _ := time.Parse(time.RFC3339, formSO.Dateso[0])
	date0 := t1.In(time.Local).Format("2006-01-02")
	t1, _ = time.Parse(time.RFC3339, formSO.Dateso[1])
	date1 := t1.In(time.Local).Format("2006-01-02") + " 23:59:59"
	// 获取数据
	var rows *sql.Rows
	if formSO.Prdno == "" {
		s := `BEGIN WITH T0_1 AS ( SELECT MF_POS.OS_DD AS [SO_DD], MF_POS.OS_NO AS [SO_NO], TF_POS.EST_ITM AS [EST_ITM_SO], TF_POS.PRD_NO, TF_POS.PRD_NAME, TF_POS.QTY AS [QTY_SO], CUST.NAME AS [CUS_NAME], TF_POS.EST_DD, TF_POS.CLS_MP_ID, MF_POS.BIL_TYPE FROM MF_POS, TF_POS, CUST WHERE MF_POS.OS_NO = TF_POS.OS_NO AND MF_POS.OS_ID = 'SO' AND MF_POS.CLS_DATE IS NOT NULL AND MF_POS.CUS_NO = CUST.CUS_NO AND ( TF_POS.EST_DD BETWEEN ? AND ? )), T0_2 AS ( SELECT DISTINCT SO_NO FROM T0_1 ), T0_3 AS ( SELECT TF_BG1.BG_FLAG, TF_BG1.OS_NO, ROW_NUMBER () OVER ( PARTITION BY TF_BG1.OS_NO ORDER BY TF_BG1.BG_DD DESC ) AS [ID] FROM MF_BG, TF_BG1, T0_2 WHERE MF_BG.BG_NO = TF_BG1.BG_NO AND MF_BG.CLS_DATE IS NOT NULL AND TF_BG1.OS_NO = T0_2.SO_NO ), T0_4 AS ( SELECT TF_POS.OS_NO, TF_POS.EST_ITM FROM T0_3, TF_POS WHERE T0_3.BG_FLAG = '4' AND T0_3.ID = 1 AND T0_3.OS_NO = TF_POS.OS_NO ), T0_7 AS ( SELECT T0_1.SO_DD, T0_1.SO_NO, T0_1.EST_ITM_SO, T0_1.PRD_NO, T0_1.PRD_NAME, T0_1.QTY_SO, T0_1.CUS_NAME, T0_1.EST_DD, T0_1.CLS_MP_ID, T0_1.BIL_TYPE FROM T0_1 LEFT JOIN T0_4 ON ( T0_1.SO_NO = T0_4.OS_NO AND T0_1.EST_ITM_SO = T0_4.EST_ITM ) WHERE T0_4.OS_NO IS NULL ), T0_5 AS ( SELECT DISTINCT PRD_NO FROM T0_7 ), T0_6 AS ( SELECT PRD_NO AS [PRD_NO], PRD_NO AS [PRD_NO1] FROM T0_5 UNION ALL ( SELECT MF_BOM.PRD_NO, TF_BOM.PRD_NO AS [PRD_NO1] FROM MF_BOM, TF_BOM, T0_5, MF_ZC WHERE MF_BOM.BOM_NO = T0_5.PRD_NO + '->' AND MF_BOM.BOM_NO = TF_BOM.BOM_NO AND ( TF_BOM.ID_NO IS NOT NULL AND TF_BOM.ID_NO <> '' ) AND TF_BOM.ID_NO = MF_ZC.BOM_NO ) ), T1 AS ( SELECT T0_7.SO_DD, T0_7.SO_NO, T0_7.EST_ITM_SO, T0_6.PRD_NO1 AS [PRD_NO], T0_7.PRD_NAME, T0_7.QTY_SO, T0_7.CUS_NAME, T0_7.EST_DD, T0_7.CLS_MP_ID, T0_7.BIL_TYPE FROM T0_7, T0_6 WHERE T0_6.PRD_NO = T0_7.PRD_NO ), T2 AS ( SELECT DISTINCT TF_MP3.MP_NO FROM T1, TF_MP3 WHERE T1.SO_NO = TF_MP3.SO_NO ), T3 AS ( SELECT MF_MO.MO_NO, MF_MO.MO_DD, MF_MO.SO_NO, MF_MO.MRP_NO AS [PRD_NO], MF_MO.QTY AS [MO_QTY], MF_MO.BAT_NO, MF_MO.BIL_NO AS [MP_NO], MF_MO.ID_NO AS [BOM_NO] FROM MF_MO, T2 WHERE MF_MO.BIL_NO = T2.MP_NO ), T4 AS ( SELECT TF_MP3.SO_NO, TF_MP3.EST_ITM, TF_MP3.PRD_NO, ROW_NUMBER () OVER ( PARTITION BY T3.MO_NO, TF_MP3.PRD_NO ORDER BY TF_POS.EST_DD, TF_POS.OS_DD, TF_POS.ITM ) AS [ID], T3.MO_NO, TF_MP3.QTY_MO, TF_POS.QTY AS [QTY_SO] FROM TF_MP3, T3, TF_POS WHERE TF_MP3.MP_NO = T3.MP_NO AND TF_MP3.PRD_NO = T3.PRD_NO AND ( TF_MP3.MO_NO = T3.MO_NO OR TF_MP3.MO_NO IS NULL OR TF_MP3.MO_NO = '' ) AND TF_MP3.SO_NO = TF_POS.OS_NO AND TF_MP3.EST_ITM = TF_POS.EST_ITM ), T5 AS ( SELECT T5_1.SO_NO, T5_1.EST_ITM, T5_1.PRD_NO, T5_1.ID AS [ID_SO_MO], T5_1.MO_NO, SUM (T5_2.QTY_SO) AS [QTY_SO_LJ] FROM T4 AS T5_1, T4 AS T5_2 WHERE ( T5_1.PRD_NO = T5_2.PRD_NO AND T5_1.MO_NO = T5_2.MO_NO AND T5_2.ID <= T5_1.ID ) GROUP BY T5_1.SO_NO, T5_1.EST_ITM, T5_1.PRD_NO, T5_1.MO_NO, T5_1.ID ) SELECT T1.SO_NO, T1.EST_ITM_SO, T1.PRD_NO, T1.PRD_NAME, T1.QTY_SO, T1.CUS_NAME, T1.EST_DD, T1.CLS_MP_ID, T5.MO_NO, T5.QTY_SO_LJ, T1.BIL_TYPE, '未完成' AS [STATUS] FROM T1 LEFT JOIN T5 ON ( T1.SO_NO = T5.SO_NO AND T1.EST_ITM_SO = T5.EST_ITM AND T1.PRD_NO = T5.PRD_NO ) WHERE ( NOT ( T1.CLS_MP_ID = 'T' AND T5.MO_NO IS NULL ) ) OR ( T1.CLS_MP_ID IS NULL AND T5.MO_NO IS NULL ) ORDER BY T1.EST_DD, T5.MO_NO END`
		rows, err = Sqlconn.Query(s, date0, date1)
	} else {
		s := `BEGIN WITH T0_1 AS ( SELECT MF_POS.OS_DD AS [SO_DD], MF_POS.OS_NO AS [SO_NO], TF_POS.EST_ITM AS [EST_ITM_SO], TF_POS.PRD_NO, TF_POS.PRD_NAME, TF_POS.QTY AS [QTY_SO], CUST.NAME AS [CUS_NAME], TF_POS.EST_DD, TF_POS.CLS_MP_ID, MF_POS.BIL_TYPE FROM MF_POS, TF_POS, CUST WHERE MF_POS.OS_NO = TF_POS.OS_NO AND MF_POS.OS_ID = 'SO' AND MF_POS.CLS_DATE IS NOT NULL AND MF_POS.CUS_NO = CUST.CUS_NO AND ( TF_POS.EST_DD BETWEEN ? AND ? ) AND TF_POS.PRD_NO = ? ), T0_2 AS ( SELECT DISTINCT SO_NO FROM T0_1 ), T0_3 AS ( SELECT TF_BG1.BG_FLAG, TF_BG1.OS_NO, ROW_NUMBER () OVER ( PARTITION BY TF_BG1.OS_NO ORDER BY TF_BG1.BG_DD DESC ) AS [ID] FROM MF_BG, TF_BG1, T0_2 WHERE MF_BG.BG_NO = TF_BG1.BG_NO AND MF_BG.CLS_DATE IS NOT NULL AND TF_BG1.OS_NO = T0_2.SO_NO ), T0_4 AS ( SELECT TF_POS.OS_NO, TF_POS.EST_ITM FROM T0_3, TF_POS WHERE T0_3.BG_FLAG = '4' AND T0_3.ID = 1 AND T0_3.OS_NO = TF_POS.OS_NO ), T0_7 AS ( SELECT T0_1.SO_DD, T0_1.SO_NO, T0_1.EST_ITM_SO, T0_1.PRD_NO, T0_1.PRD_NAME, T0_1.QTY_SO, T0_1.CUS_NAME, T0_1.EST_DD, T0_1.CLS_MP_ID, T0_1.BIL_TYPE FROM T0_1 LEFT JOIN T0_4 ON ( T0_1.SO_NO = T0_4.OS_NO AND T0_1.EST_ITM_SO = T0_4.EST_ITM ) WHERE T0_4.OS_NO IS NULL ), T0_5 AS ( SELECT DISTINCT PRD_NO FROM T0_7 ), T0_6 AS ( SELECT PRD_NO AS [PRD_NO], PRD_NO AS [PRD_NO1] FROM T0_5 UNION ALL ( SELECT MF_BOM.PRD_NO, TF_BOM.PRD_NO AS [PRD_NO1] FROM MF_BOM, TF_BOM, T0_5, MF_ZC WHERE MF_BOM.BOM_NO = T0_5.PRD_NO + '->' AND MF_BOM.BOM_NO = TF_BOM.BOM_NO AND ( TF_BOM.ID_NO IS NOT NULL AND TF_BOM.ID_NO <> '' ) AND TF_BOM.ID_NO = MF_ZC.BOM_NO ) ), T1 AS ( SELECT T0_7.SO_DD, T0_7.SO_NO, T0_7.EST_ITM_SO, T0_6.PRD_NO1 AS [PRD_NO], T0_7.PRD_NAME, T0_7.QTY_SO, T0_7.CUS_NAME, T0_7.EST_DD, T0_7.CLS_MP_ID, T0_7.BIL_TYPE FROM T0_7, T0_6 WHERE T0_6.PRD_NO = T0_7.PRD_NO ), T2 AS ( SELECT DISTINCT TF_MP3.MP_NO FROM T1, TF_MP3 WHERE T1.SO_NO = TF_MP3.SO_NO ), T3 AS ( SELECT MF_MO.MO_NO, MF_MO.MO_DD, MF_MO.SO_NO, MF_MO.MRP_NO AS [PRD_NO], MF_MO.QTY AS [MO_QTY], MF_MO.BAT_NO, MF_MO.BIL_NO AS [MP_NO], MF_MO.ID_NO AS [BOM_NO] FROM MF_MO, T2 WHERE MF_MO.BIL_NO = T2.MP_NO ), T4 AS ( SELECT TF_MP3.SO_NO, TF_MP3.EST_ITM, TF_MP3.PRD_NO, ROW_NUMBER () OVER ( PARTITION BY T3.MO_NO, TF_MP3.PRD_NO ORDER BY TF_POS.EST_DD, TF_POS.OS_DD, TF_POS.ITM ) AS [ID], T3.MO_NO, TF_MP3.QTY_MO, TF_POS.QTY AS [QTY_SO] FROM TF_MP3, T3, TF_POS WHERE TF_MP3.MP_NO = T3.MP_NO AND TF_MP3.PRD_NO = T3.PRD_NO AND ( TF_MP3.MO_NO = T3.MO_NO OR TF_MP3.MO_NO IS NULL OR TF_MP3.MO_NO = '' ) AND TF_MP3.SO_NO = TF_POS.OS_NO AND TF_MP3.EST_ITM = TF_POS.EST_ITM ), T5 AS ( SELECT T5_1.SO_NO, T5_1.EST_ITM, T5_1.PRD_NO, T5_1.ID AS [ID_SO_MO], T5_1.MO_NO, SUM (T5_2.QTY_SO) AS [QTY_SO_LJ] FROM T4 AS T5_1, T4 AS T5_2 WHERE ( T5_1.PRD_NO = T5_2.PRD_NO AND T5_1.MO_NO = T5_2.MO_NO AND T5_2.ID <= T5_1.ID ) GROUP BY T5_1.SO_NO, T5_1.EST_ITM, T5_1.PRD_NO, T5_1.MO_NO, T5_1.ID ) SELECT T1.SO_NO, T1.EST_ITM_SO, T1.PRD_NO, T1.PRD_NAME, T1.QTY_SO, T1.CUS_NAME, T1.EST_DD, T1.CLS_MP_ID, T5.MO_NO, T5.QTY_SO_LJ, T1.BIL_TYPE, '未完成' AS [STATUS] FROM T1 LEFT JOIN T5 ON ( T1.SO_NO = T5.SO_NO AND T1.EST_ITM_SO = T5.EST_ITM AND T1.PRD_NO = T5.PRD_NO ) WHERE ( NOT ( T1.CLS_MP_ID = 'T' AND T5.MO_NO IS NULL ) ) OR ( T1.CLS_MP_ID IS NULL AND T5.MO_NO IS NULL ) ORDER BY T1.EST_DD, T5.MO_NO END`
		rows, err = Sqlconn.Query(s, date0, date1, formSO.Prdno)
	}
	if err != nil {
		log.Fatal("GetSoTableJsonFromDB rows err: ", err.Error())
	}
	defer rows.Close()
	// 取得所有mono
	var monomap = make(map[string]int)
	for rows.Next() {
		var sono, prdno, prdname, cusname, estdd, clsmpid, mono, biltype, status string
		var estitmso int
		var qtyso, qtysolj float32
		rows.Scan(&sono, &estitmso, &prdno, &prdname, &qtyso, &cusname, &estdd, &clsmpid, &mono, &qtysolj, &biltype, &status)
		// 繁简转换
		cusname, _ = tw2s.Convert(cusname)
		prdname, _ = tw2s.Convert(prdname)
		// 设置订单类型
		switch biltype {
		case "01":
			biltype = "正常订单"
		case "02":
			biltype = "样品订单"
		}
		// 取mono
		if mono != "" {
			monomap["TF_WR.MO_NO = '" + mono + "'"]+=1
		}
		sodatas = append(sodatas, JsSoData{sono, estitmso, prdno, prdname, qtyso, cusname, estdd[:10], clsmpid, mono, qtysolj, biltype, status})
	}
	if len(monomap) > 0 {
		// 查询所有指令单包装上帐数量
		monolist := "("
		for k, _ := range monomap {
			monolist += k+" OR "
		}
		l := len(monolist) - 4
		s1 := fmt.Sprintf("SELECT TF_WR.MO_NO, SUM(TF_WR.QTY_FIN) AS [QTY_BZ] FROM TF_WR, TF_ZC WHERE TF_WR.ID_NO = TF_ZC.BOM_NO AND TF_WR.ZC_NO = TF_ZC.ZC_NO AND (TF_ZC.ZC_NO_DN IS NULL OR TF_ZC.ZC_NO_DN = '') AND %s GROUP BY TF_WR.MO_NO", monolist[:l]+ ")")
		rows1, err := Sqlconn.Query(s1)
		if err != nil {
			log.Fatal("GetSoTableJsonFromDB in get bz rows err: ", err.Error())
		}
		defer rows1.Close()
		for rows1.Next() {
			var mo string
			var bz float32
			rows1.Scan(&mo, &bz)
			for i, v := range sodatas {
				if v.Mono == mo && v.Mono != "" {
					if bz < v.QtySoLj {
						sodatas[i].Status = "未完成"
					} else {
						sodatas[i].Status = "已完成"
					}
				}
			}
		}
	}
	// sodatas = []JsSoData{{"Q-001"}, {"Q-002"}, {"Q-003"}, {"Q-004"}, {"Q-005"}, {"Q-006"}, {"Q-007"}, {"Q-008"}, {"Q-009"}, {"Q-010"}, {"Q-011"}, {"Q-012"}, {"Q-013"}, {"Q-014"}, {"Q-015"}, {"Q-016"}, {"Q-017"}, {"Q-018"}}
	// log.Println(sodatas)
	return sodatas
}

func GetTzTableFromDB(mono, prdno string) ([]JsTzData, []JsBcpSData, float32) {
	var tzdatas []JsTzData
	var bcpsdatas []JsBcpSData
	var qtycp float32
	// 获取通知单明细
	s := `BEGIN WITH T0 AS ( SELECT MF_TZ.TZ_NO, MF_TZ.ZC_ITM, DEPT.NAME AS [DEP_NAME], ZC_NO.NAME AS [ZC_NAME], ISNULL(MF_TZ.QTY, 0) AS [QTY], ISNULL(MF_TZ.QTY_FIN, 0) AS [QTY_FIN], ISNULL(MF_TZ.QTY_LOST, 0) AS [QTY_LOST], MF_TZ.MRP_NO AS [PRD_NO], MF_TZ.ZC_NO FROM MF_TZ, ZC_NO, DEPT WHERE MF_TZ.MO_NO = ? AND ( MF_TZ.BIL_ID <> 'TR' OR MF_TZ.BIL_ID IS NULL ) AND MF_TZ.ZC_NO = ZC_NO.ZC_NO AND MF_TZ.DEP = DEPT.DEP ), T1 AS ( SELECT T0.TZ_NO, T0.ZC_ITM, T0.DEP_NAME, T0.ZC_NAME, T0.QTY, T0.QTY_FIN, T0.QTY_LOST, TT1.MY_DINGE FROM T0 LEFT JOIN (SELECT UP_DEF.PRD_NO, UP_DEF.BZ_KND AS [ZC_NO], ISNULL(UP_DEF_Z.MY_DINGE, 0) AS [MY_DINGE], ROW_NUMBER () OVER ( PARTITION BY UP_DEF_Z.PRICE_ID, UP_DEF_Z.CUS_NO, UP_DEF_Z.CUR_ID, UP_DEF_Z.PRD_NO, UP_DEF_Z.BZ_KND, UP_DEF_Z.QTY, UP_DEF_Z.BIL_TYPE ORDER BY UP_DEF_Z.S_DD DESC ) AS ID FROM UP_DEF, UP_DEF_Z WHERE UP_DEF.CHK_MAN IS NOT NULL AND UP_DEF.PRICE_ID = '3' AND UP_DEF.CUS_NO = '0000' AND (UP_DEF.BIL_TYPE = '01') AND UP_DEF.QTY = 1 AND UP_DEF.PRICE_ID = UP_DEF_Z.PRICE_ID AND UP_DEF.CUR_ID = UP_DEF_Z.CUR_ID AND UP_DEF.CUS_NO = UP_DEF_Z.CUS_NO AND UP_DEF.PRD_NO = UP_DEF_Z.PRD_NO AND UP_DEF.BZ_KND = UP_DEF_Z.BZ_KND AND UP_DEF.QTY = UP_DEF_Z.QTY AND UP_DEF.S_DD = UP_DEF_Z.S_DD AND UP_DEF.BIL_TYPE = UP_DEF_Z.BIL_TYPE AND UP_DEF.PRD_NO = ? ) AS TT1 ON ( TT1.ID = 1 AND T0.PRD_NO = TT1.PRD_NO AND T0.ZC_NO = TT1.ZC_NO) ), T2 AS ( SELECT TF_SCPG.BIL_NO, ISNULL(SUM(TF_SCPG.QTY), 0) AS [QTY_PGS] FROM TF_SCPG, T1 WHERE TF_SCPG.BIL_NO = T1.TZ_NO GROUP BY TF_SCPG.BIL_NO ), T3 AS ( SELECT TZERR.TZ_NO, ISNULL(SUM(TF_IJ.QTY), 0) AS [QTY_BF] FROM MF_IJ, TF_IJ, TZERR, T1 WHERE TZERR.TZ_NO = T1.TZ_NO AND MF_IJ.BIL_NO = TZERR.TR_NO AND MF_IJ.IJ_ID = 'XF' AND TF_IJ.IJ_ID = 'XF' AND MF_IJ.IJ_NO = TF_IJ.IJ_NO GROUP BY TZERR.TZ_NO ) SELECT T1.TZ_NO, T1.DEP_NAME, T1.ZC_NAME, T1.QTY, T1.QTY_FIN, T1.QTY_LOST, ISNULL(T3.QTY_BF, 0) AS [QTY_BF], ( T1.QTY - ISNULL(T1.QTY_FIN, 0) - ISNULL(T3.QTY_BF, 0) ) AS [QTY_SY], ISNULL(T2.QTY_PGS, 0) AS [QTY_PGS], T1.MY_DINGE FROM T1 FULL JOIN T2 ON T1.TZ_NO = T2.BIL_NO FULL JOIN T3 ON T1.TZ_NO = T3.TZ_NO ORDER BY T1.ZC_ITM END`
	rows, err := Sqlconn.Query(s, mono, prdno)
	if err != nil {
		log.Fatal("GetTzTableFromDB tz rows err: ", err.Error())
	}
	defer rows.Close()
	// 处理数据
	for rows.Next() {
		var tzno, depname, zcname string
		var qty, qtyfin, qtylost, qtybf, qtysy, qtypgs, mydinge float32
		rows.Scan(&tzno, &depname, &zcname, &qty, &qtyfin, &qtylost, &qtybf, &qtysy, &qtypgs, &mydinge)
		// 繁简转换
		depname, _ = tw2s.Convert(depname)
		zcname, _ = tw2s.Convert(zcname)
		tzdatas = append(tzdatas, JsTzData{tzno, depname, zcname, qty, qtyfin, qtylost, qtybf, qtysy, qtypgs, mydinge})
	}
	// 获取半成品各序库存总数
	s1 := `SELECT BAT_REC1.PRD_NO, PRDT.SPC, SUM(ISNULL(BAT_REC1.QTY_IN, 0) - ISNULL(BAT_REC1.QTY_OUT, 0)) AS [QTY] FROM BAT_REC1, PRDT WHERE BAT_REC1.PRD_NO LIKE ? + '%' AND BAT_REC1.PRD_NO = PRDT.PRD_NO AND (BAT_REC1.WH LIKE '11%' OR BAT_REC1.WH LIKE '12%' OR BAT_REC1.WH LIKE '13%') AND (BAT_REC1.WH <> '11  999999' AND BAT_REC1.WH <> '12  999999' AND BAT_REC1.WH <> '13  999999') AND (ISNULL(BAT_REC1.QTY_IN, 0) - ISNULL(BAT_REC1.QTY_OUT, 0) > 0) GROUP BY BAT_REC1.PRD_NO, PRDT.SPC`
	rows1, err := Sqlconn.Query(s1, prdno)
	if err != nil {
		log.Fatal("GetTzTableFromDB bcp rows err: ", err.Error())
	}
	defer rows1.Close()
	for rows1.Next() {
		var prdno, zcname string
		var qty float32
		rows1.Scan(&prdno, &zcname, &qty)
		// 繁简转换
		zcname, _ = tw2s.Convert(zcname)	
		bcpsdatas = append(bcpsdatas, JsBcpSData{prdno, zcname, qty})
	}
	// 获取成品库库存总数
	s2 := `SELECT ISNULL(SUM(ISNULL(QTY_IN, 0) - ISNULL(QTY_OUT, 0)), 0) AS [QTY] FROM DB_CH.dbo.BAT_REC1 WHERE PRD_NO = ? AND (WH = '02' OR WH LIKE '02  %') AND (WH <> '02  000000') AND ISNULL(QTY_IN, 0) - ISNULL(QTY_OUT, 0) > 0`
	err = Sqlconn.QueryRow(s2, prdno).Scan(&qtycp)
	if err != nil {
		log.Fatal("GetTzTableFromDB cp rows err: ", err.Error())
	}
	return tzdatas, bcpsdatas, qtycp
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
