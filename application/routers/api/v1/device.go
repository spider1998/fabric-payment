package v1

import (
	bc "application/blockchain"
	"application/lib"
	"application/pkg/app"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"

	"golang.org/x/text/transform"

	"github.com/gin-gonic/gin"
)

var WeekDays = map[string]string{
	"Monday":    "1",
	"Tuesday":   "2",
	"Wednesday": "3",
	"Thursday":  "4",
	"Friday":    "5",
	"Saturday":  "6",
	"Sunday":    "0",
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func padUnMarshal(str string, v interface{}) (err error) {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	list := strings.Split(str, ",")
	if len(list) == 0 {
		return
	}
	e := reflect.Indirect(reflect.ValueOf(v))
	for _, li := range list {
		res := regexp.MustCompile(`\".*?\"`)
		kv := res.FindAllString(li, -1)
		if len(kv) != 2 {
			break
		}
		k := strings.Replace(kv[0], "\"", "", -1)
		fmt.Println(k)
		v := strings.Replace(kv[1], "\"", "", -1)
		fmt.Println(v)
		e.FieldByName(k).SetString(v)
	}
	return
}

func encryption(cardNo string) string {
	cardKey := string(smsEncDec([]byte(cardNo), []byte("shengdian")))
	h := md5.New()
	h.Write([]byte(cardKey)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

func smsEncDec(data, key []byte) []byte {
	enc := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		enc[i] = data[i] ^ key[i%len(key)]
	}
	return enc
}

//综合消费接口
func ConsumTransactions(c *gin.Context) {
	appG := app.Gin{C: c}
	machineCode := c.Request.Header.Get("Client-ID")
	var req lib.ConsumTransactionsCond
	var res lib.ConsumTransactionsResponse
	s, _ := ioutil.ReadAll(c.Request.Body)
	str := string(s)
	fmt.Println(str)
	err := padUnMarshal(str, &req)
	if err != nil {
		return
	}
	if req.CardNo != "" {
		req.CardNo = encryption(req.CardNo)
	}
	req.MachineCode = machineCode
	var b []byte
	b, err = json.Marshal(req)
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体出错%s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, b)

	resp, err := bc.ChannelExecute("consumTransactions", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &res); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	b, err = json.MarshalIndent(&res, "", " ")
	if err != nil {
		return
	}
	bs, _ := Utf8ToGbk(b)
	c.Header("Content-Type", "application/json; charset=gb2312")

	c.Writer.Write(bs)
}

//获取系统时间
func GetTimes(c *gin.Context) {
	//c.Response.Header().Set("Content-Type","application/json; charset=gb2312")
	c.Request.Header.Set("Content-Type", "application/json")
	c.Header("Content-Type", "application/json; charset=gb2312")
	c.Header("Client-ID", "1")
	var res lib.GetTimesResponse
	res.Status = "100"
	res.Time = time.Now().Format("20060102150405") + WeekDays[time.Now().Weekday().String()]

	b, err := json.MarshalIndent(&res, "", " ")
	if err != nil {
		//sendError(w, http.StatusBadRequest, CodeErrorServer, "数据转换出错")
		return
	}
	c.Writer.Write(b)
}
