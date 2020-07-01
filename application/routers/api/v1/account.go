/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/5 12:55 上午
 * @Description: 账户信息相关接口
 */
package v1

import (
	bc "application/blockchain"
	"application/lib"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountIdBody struct {
	AccountId string `json:"accountId"`
}

type AccountRequestBody struct {
	Args []AccountIdBody `json:"args"`
}

// @Summary 获取账户信息
// @Param account body AccountRequestBody true "account"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/queryAccountList [post]
func QueryAccountList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(AccountRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	for _, val := range body.Args {
		bodyBytes = append(bodyBytes, []byte(val.AccountId))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryAccountList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryAccountsList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(AccountRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	for _, val := range body.Args {
		bodyBytes = append(bodyBytes, []byte(val.AccountId))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryAccountsList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func Recharge(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(lib.Tem)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var b []byte
	b, err := json.Marshal(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体出错%s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, b)

	resp, err := bc.ChannelExecute("recharge", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func UpdateAccounts(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(lib.UpdateAccountsRequest)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var b []byte
	b, err := json.Marshal(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体出错%s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, b)

	resp, err := bc.ChannelExecute("updateAccount", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func CreateAccounts(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(lib.Account)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var b []byte
	b, err := json.Marshal(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("解析请求体出错%s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, b)

	resp, err := bc.ChannelExecute("createAccount", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
