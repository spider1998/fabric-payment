package routers

import (
	"chaincode/lib"
	"chaincode/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	uuid "github.com/satori/go.uuid"
)

func ConsumTransactions(stub shim.ChaincodeStubInterface, body [][]byte) pb.Response {
	var req lib.ConsumTransactionsCond
	var res lib.ConsumTransactionsResponse
	err := json.Unmarshal(body[1], &req)
	if err != nil {
		return shim.Error(fmt.Sprintf("消费机-反序列化出错: %s", err))
	}
	if req.PayType == "2" {
		fmt.Println(req.QR)
		fmt.Println("二维码消费")
		res = qrConsumption(stub, req, req.MachineCode)
	} else if req.PayType == "0" {
		if req.Mode == "0" {
			fmt.Println("刷卡消费")
			res = cardConsumption(stub, req, req.MachineCode)
		} else if req.Mode == "3" {
			fmt.Println("信息查询")
			res = cardInfo(stub, req)
		}
	} else if req.PayType == "3" {
		fmt.Println("二维码消费确认")
		res = qrConsumptionConfirm(stub, req, req.MachineCode)
	}
	result, err := json.Marshal(res)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	return shim.Success(result)
}

func qrConsumption(stub shim.ChaincodeStubInterface, req lib.ConsumTransactionsCond, machineCode string) (res lib.ConsumTransactionsResponse) {
	res.Status = "101"
	res.Msg = "暂不支持"
	return
}
func qrConsumptionConfirm(stub shim.ChaincodeStubInterface, req lib.ConsumTransactionsCond, machineCode string) (res lib.ConsumTransactionsResponse) {
	res.Status = "101"
	res.Msg = "暂不支持"
	return
}

//信息查询
func cardInfo(stub shim.ChaincodeStubInterface, req lib.ConsumTransactionsCond) (res lib.ConsumTransactionsResponse) {
	//查询卡状态是否正常
	var account lib.Accounts
	queryString := fmt.Sprintf("{\"selector\":{\"card_id\":\"%s\",\"table\":\"%s\"}}", req.CardNo, "accounts")
	// 查询数据
	result, err := getAccountsByQueryString(stub, queryString)
	if err != nil {
		res.Status = "101"
		res.Msg = "卡查询出错"
		return
	}
	if result == nil {
		res.Status = "101"
		res.Msg = "账户不存在"
		return
	}
	err = json.Unmarshal(result, &account)
	if err != nil {
		res.Status = "101"
		res.Msg = "解析错误"
		return
	}
	if account.Customer.Cardstatus != "1" {
		res.Status = "101"
		res.Msg = "卡状态错误"
		return
	}
	if account.Status != "1" {
		res.Status = "101"
		res.Msg = "账户不正常"
		return
	}
	res = lib.ConsumTransactionsResponse{
		StartTime:    "20191105163029",
		EndTime:      "20991231000000",
		Status:       "100",
		Msg:          "",
		Name:         account.Name,
		Money:        strconv.FormatFloat(float64(account.Balance), 'f', -1, 64),
		Subsidy:      "0.00",
		Times:        "0",
		Integral:     "0",
		Dept:         account.Customer.Department,
		Discount:     "100",
		QR:           "",
		CardType:     "",
		InTime:       "",
		OutTime:      "",
		Amount:       "0",
		QROrder:      "",
		PayResults:   "",
		Level:        "0",
		CardNumber:   "",
		QRCodeNumber: "",
		VoiceID:      "",
		Text:         "查询成功",
	}
	return
}

//刷卡消费
func cardConsumption(stub shim.ChaincodeStubInterface, req lib.ConsumTransactionsCond, machineCode string) (res lib.ConsumTransactionsResponse) {
	//查询卡状态是否正常
	var account lib.Accounts
	queryString := fmt.Sprintf("{\"selector\":{\"owner.cardID\":\"%s\",\"table\":\"%s\"}}", req.CardNo, "accounts")
	// 查询数据
	result, err := getAccountsByQueryString(stub, queryString)
	if err != nil {
		res.Status = "101"
		res.Msg = "卡查询出错"
		return
	}
	if result == nil {
		res.Status = "101"
		res.Msg = "账户不存在"
		return
	}
	err = json.Unmarshal(result, &account)
	if err != nil {
		res.Status = "101"
		res.Msg = "解析错误"
		return
	}
	if account.Customer.Cardstatus != "1" {
		res.Status = "101"
		res.Msg = "卡状态错误"
		return
	}
	if account.Status != "1" {
		res.Status = "101"
		res.Msg = "账户不正常"
		return
	}
	//比对余额
	fAmount, _ := strconv.ParseFloat(req.Amount, 64)
	amount := int(fAmount * 100)
	if amount > account.Balance {
		res.Status = "101"
		res.Msg = "余额不足"
		return
	}
	//todo 检查消费机是否连接正常
	//生成交易记录
	account.Balance -= amount
	dealTime := time.Now().Unix()
	timeNow := time.Unix(dealTime, 0).Format("2006-01-02 15:04:05")
	rc := lib.Record{
		ID:            uuid.NewV1().String(),
		AccountID:     account.ID,
		TransactionID: req.Order,
		JobNumber:     account.EmployeeID,
		CardNumber:    account.Customer.CardID,
		CustomerID:    account.CustomerID,
		CustomerName:  account.Customer.Name,
		Type:          lib.TypeTransaction,
		Amounts:       amount,
		Status:        lib.StatusRefund,
		Operator:      lib.OperatorSys,
		RealMoney:     amount,
		Balance:       account.Balance,
		MerchantID:    "",
		DeviceID:      "",
		DeviceCode:    machineCode,
		DealTime:      dealTime,

		Remark: "刷卡交易",
	}
	rc.ShowDate = timeNow[0:10]
	rc.UpdateTime = timeNow
	if err := utils.WriteLedger(rc, stub, lib.RecordKey, []string{rc.ID}); err != nil {
		res.Status = "101"
		res.Msg = "改账户失败"
		return
	}

	//修改账户余额
	account.UpdateTime = timeNow
	if err := utils.WriteLedger(account, stub, lib.AccountsKey, []string{account.ID}); err != nil {
		res.Status = "101"
		res.Msg = "改余额失败"
		return
	}

	//todo 修改消费机统计

	res = lib.ConsumTransactionsResponse{
		StartTime:    "20191105163029",
		EndTime:      "20991231000000",
		Status:       "100",
		Msg:          "",
		Name:         account.Name,
		Money:        strconv.FormatFloat(float64(account.Balance), 'f', -1, 64),
		Subsidy:      "0.00",
		Times:        "0",
		Integral:     "0",
		Dept:         account.Customer.Department,
		Discount:     "100",
		QR:           "",
		CardType:     "",
		InTime:       "",
		OutTime:      "",
		Amount:       req.Amount,
		QROrder:      "",
		PayResults:   "",
		Level:        "0",
		CardNumber:   "",
		QRCodeNumber: "",
		VoiceID:      "",
		Text:         "支付成功",
	}
	return
}
