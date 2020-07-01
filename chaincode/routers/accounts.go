package routers

import (
	"bytes"
	"chaincode/lib"
	"chaincode/utils"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	//uuid "github.com/satori/go.uuid"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getAccountsByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

func checkParams(cash, virtual int) bool {
	if cash < 0 {
		return false
	}
	if virtual < 0 {
		return false
	}
	if cash == 0 && virtual == 0 {
		return false
	}
	return true
}

func Recharge(stub shim.ChaincodeStubInterface, body [][]byte) pb.Response {
	t := []lib.Tem{}
	err := json.Unmarshal(body[1], &t)
	if err != nil {
		return shim.Error(fmt.Sprintf("账户充值-反序列化出错: %s", err))
	}
	timeStamp := time.Now().Unix()
	timeNow := time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
	for _, v := range t {
		if !checkParams(v.Cash, v.Virtual) {
			return shim.Error(fmt.Sprintf("携带金额参数出错"))
		}
		var account lib.Accounts
		queryString := fmt.Sprintf("{\"selector\":{\"id\":\"%s\",\"table\":\"%s\"}}", v.ID, "accounts")
		// 查询数据
		result, err := getAccountsByQueryString(stub, queryString)
		if err != nil {
			return shim.Error("根据id查询账户信息时发生错误" + err.Error())
		}
		if result == nil {
			return shim.Error("修改账户-人员不存在")
		}
		err = json.Unmarshal(result, &account)
		if err != nil {
			return shim.Error("修改前账户查询数据解析失败" + err.Error())
		}
		preBalance := account.Balance
		balance := account.Balance + v.Cash + v.Virtual
		account.Balance = balance
		//有补助
		if v.Virtual > 0 {
			account.Virtual = v.Virtual
			account.VirtualTime = timeNow
		}
		if err := utils.WriteLedger(account, stub, lib.AccountsKey, []string{account.ID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}

		showTime := timeNow
		showDate := ""
		if showTime != "" {
			showDate = showTime[0:10]
		}
		preBalance += v.Cash
		if v.Cash > 0 { //现金充值
			record := lib.Record{
				ID:           uuid.NewV1().String(),
				AccountID:    account.ID,
				CustomerID:   account.CustomerID,
				CustomerName: account.Customer.Name,
				Department:   account.Department,
				JobNumber:    account.EmployeeID,
				VirtualMoney: -v.Virtual,
				RealMoney:    -v.Cash,
				Amounts:      -(v.Cash),
				DealTime:     timeStamp,
				Balance:      preBalance,
				Operator:     lib.OperatorSys,
				Type:         lib.TypeRecharge,
				Status:       lib.StatusRefund,
				CardNumber:   account.Customer.CardID,
				ShowTime:     showTime,
				ShowDate:     showDate,
				Remark:       "账户现金充值",
				UpdateTime:   timeNow,
			}
			record.ShowDate = time.Now().Format("2006-01-02 15:04:05")[0:10]
			record.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			if err := utils.WriteLedger(record, stub, lib.RecordKey, []string{record.ID}); err != nil {
				return shim.Error(fmt.Sprintf("%s", err))
			}
		}
		preBalance += v.Virtual
		if v.Virtual > 0 { //补助
			record := lib.Record{
				ID:           uuid.NewV1().String(),
				AccountID:    account.ID,
				CustomerID:   account.CustomerID,
				CustomerName: account.Customer.Name,
				Department:   account.Department,
				JobNumber:    account.EmployeeID,
				VirtualMoney: -v.Virtual,
				RealMoney:    -v.Cash,
				Amounts:      -(v.Virtual),
				Balance:      preBalance,
				DealTime:     timeStamp,
				Operator:     lib.OperatorSys,
				Type:         lib.TypeSubsidy,
				Status:       lib.StatusRefund,
				CardNumber:   account.Customer.CardID,
				Remark:       "系统补助",
				ShowTime:     showTime,
				ShowDate:     showDate,
				UpdateTime:   timeNow,
			}
			record.ShowDate = time.Now().Format("2006-01-02 15:04:05")[0:10]
			record.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			if err := utils.WriteLedger(record, stub, lib.RecordKey, []string{record.ID}); err != nil {
				return shim.Error(fmt.Sprintf("%s", err))
			}
		}
	}
	return shim.Success([]byte("OK"))
}

func UpdateAccount(stub shim.ChaincodeStubInterface, body [][]byte) pb.Response {
	var account lib.UpdateAccountsRequest
	err := json.Unmarshal(body[1], &account)
	if err != nil {
		return shim.Error(fmt.Sprintf("修改账户-反序列化出错: %s", err))
	}
	var preAccount lib.Accounts
	queryString := fmt.Sprintf("{\"selector\":{\"id\":\"%s\",\"table\":\"%s\"}}", account.ID, "accounts")
	// 查询数据
	result, err := getAccountsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据id查询账户信息时发生错误" + err.Error())
	}
	if result == nil {
		return shim.Error("修改账户-人员不存在")
	}
	err = json.Unmarshal(result, &preAccount)
	if err != nil {
		return shim.Error("修改前账户查询数据解析失败" + err.Error())
	}
	//部门有变动
	if account.Department != preAccount.Department {
		preAccount.Department = account.Department
		preAccount.Customer.Department = account.Department
	}
	dealTime := time.Now().Unix()
	if account.Name != preAccount.Name {
		preAccount.Name = account.Name
		preAccount.Customer.Name = account.Name
	}
	//身份证号变动
	if account.IDCard != "" {
		preAccount.Customer.IDCard = account.IDCard
	}
	//手机号有变动
	if account.Cellphone != "" && account.Cellphone != preAccount.Customer.Cellchone {
		queryString := fmt.Sprintf("{\"selector\":{\"owner.cellphone\":\"%s\",\"id\":{\"$ne\":\"%s\"},\"table\":\"%s\"}}", account.Cellphone, account.ID, "accounts")
		// 查询数据
		result, err := getAccountsByQueryString(stub, queryString)
		if err != nil {
			return shim.Error("根据手机号+id查询账户信息时发生错误" + err.Error())
		}
		if result != nil {
			return shim.Error("修改账户-手机号已存在")
		}
		preAccount.Customer.Cellchone = account.Cellphone
	}
	//人员编号有变动
	if account.EmployeeID != "" {
		preAccount.Customer.EmployeeID = account.EmployeeID
		preAccount.EmployeeID = account.EmployeeID
	}
	//账户状态（挂失，解挂，注销）
	if account.Status != "" {
		if account.Status == "3" { //注销
			preAccount.Customer.Cardstatus = "4"
		}
		preAccount.Status = account.Status
	}
	preAccount.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	//卡费
	if account.CardFee != 0 {
		preAccount.CardFee = account.CardFee
		preAccount.Customer.CardFee = account.CardFee
	}
	//更新账户
	if err := utils.WriteLedger(preAccount, stub, lib.AccountsKey, []string{account.ID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if account.CardFee != 0 {
		rc := lib.Record{
			ID:            uuid.NewV1().String(),
			AccountID:     account.ID,
			TransactionID: uuid.NewV4().String(),
			JobNumber:     account.EmployeeID,
			CardNumber:    preAccount.Customer.CardID,
			CustomerID:    preAccount.CustomerID,
			CustomerName:  preAccount.Customer.Name,
			Type:          lib.Typecard,
			Amounts:       account.CardFee,
			Status:        lib.StatusRefund,
			Operator:      lib.OperatorSys,
			RealMoney:     account.CardFee,
			Balance:       account.Balance,
			DealTime:      dealTime,
			Remark:        "卡费",
		}
		rc.ShowDate = time.Now().Format("2006-01-02 15:04:05")[0:10]
		rc.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
		if err := utils.WriteLedger(rc, stub, lib.RecordKey, []string{rc.ID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}
	accounts, err := json.Marshal(preAccount)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	return shim.Success(accounts)
}

//新建消费账户
func CreateAccount(stub shim.ChaincodeStubInterface, body [][]byte) pb.Response {
	var account lib.Accounts
	err := json.Unmarshal(body[1], &account)
	if err != nil {
		return shim.Error(fmt.Sprintf("创建账户-反序列化出错: %s", err))
	}
	if account.Customer.Name == "" {
		return shim.Error(fmt.Sprintf("创建账户-账户名称不能为空: %s", err))
	}
	if account.Customer.EmployeeID == "" {
		return shim.Error(fmt.Sprintf("创建账户-人员编号不能为空: %s", err))
	}
	if account.Customer.CardID == "" {
		return shim.Error(fmt.Sprintf("创建账户-人员卡号不能为空: %s", err))
	}
	if account.Department == "" {
		return shim.Error(fmt.Sprintf("创建账户-人员部门不能为空: %s", err))
	}
	queryString := fmt.Sprintf("{\"selector\":{\"owner.cardID\":\"%s\",\"table\":\"%s\"}}", account.Customer.CardID, "accounts")

	// 查询数据
	result, err := getAccountsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据卡号查询账户信息时发生错误" + err.Error())
	}
	if result != nil {
		return shim.Error("创建账户-人员卡号已存在")
	}
	queryCellphone := fmt.Sprintf("{\"selector\":{\"owner.cellphone\":\"%s\",\"owner.cardstatus\":\"%s\",\"table\":\"%s\"}}", account.Customer.Cellchone, "1", "accounts")

	// 查询数据
	result, err = getAccountsByQueryString(stub, queryCellphone)
	if err != nil {
		return shim.Error("根据手机号查询账户信息时发生错误" + err.Error())
	}
	if result != nil {
		return shim.Error("创建账户-手机号已存在")
	}
	if account.Status == "" {
		account.Status = "1"
	}
	account.ID = uuid.NewV1().String()
	account.CustomerID = uuid.NewV1().String()
	account.Customer.ID = account.CustomerID
	account.Customer.AccountID = account.ID
	account.Customer.Department = account.Department
	account.Customer.CardFee = account.CardFee
	account.Customer.Types = 1
	account.Table = "accounts"

	if account.Customer.CardNumber != "" {
		account.Customer.CardID = account.Customer.CardNumber
		account.Customer.CardCreateTime = time.Now().Format("2006-01-02 15:04:05")
		account.Customer.Cardstatus = "1"
	}
	account.Name = account.Customer.Name
	account.EmployeeID = account.Customer.EmployeeID
	account.EmployeeStatus = 1
	account.CreateTime = time.Now().Format("2006-01-02 :15:04:05")
	account.UpdateTime = time.Now().Format("2006-01-02 :15:04:05")
	//TODO 账户部门是否存在
	dealTime := time.Now().Unix()
	timeNow := time.Unix(dealTime, 0).Format("2006-01-02 15:04:05")
	//卡费不为空
	if account.CardFee != 0 {
		rc := lib.Record{
			ID:            uuid.NewV1().String(),
			AccountID:     account.ID,
			TransactionID: uuid.NewV4().String(),
			JobNumber:     account.EmployeeID,
			CardNumber:    account.Customer.CardID,
			CustomerID:    account.CustomerID,
			CustomerName:  account.Customer.Name,
			Type:          lib.Typecard,
			Amounts:       account.CardFee,
			Status:        lib.StatusRefund,
			Operator:      lib.OperatorSys,
			RealMoney:     account.CardFee,
			Balance:       account.Balance,
			DealTime:      dealTime,
			Remark:        "开卡卡费",
		}
		rc.ShowDate = timeNow[0:10]
		rc.UpdateTime = timeNow
		if err := utils.WriteLedger(rc, stub, lib.RecordKey, []string{rc.ID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}
	account.Balance += account.Cash
	//记录充值
	if account.Cash != 0 {
		rc := lib.Record{
			ID:            uuid.NewV1().String(),
			AccountID:     account.ID,
			TransactionID: uuid.NewV4().String(),
			JobNumber:     account.EmployeeID,
			CardNumber:    account.Customer.CardID,
			CustomerID:    account.CustomerID,
			CustomerName:  account.Customer.Name,
			Type:          lib.TypeRecharge,
			Amounts:       -account.Cash,
			Status:        lib.StatusRefund,
			Operator:      lib.OperatorSys,
			RealMoney:     account.Cash,
			Balance:       account.Balance,
			DealTime:      dealTime,
			Remark:        "开卡现金充值",
		}
		rc.ShowDate = timeNow[0:10]
		rc.UpdateTime = timeNow
		if err := utils.WriteLedger(rc, stub, lib.RecordKey, []string{rc.ID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}
	account.Balance += account.Virtual
	//记录补助
	if account.Virtual != 0 {
		rc := lib.Record{
			ID:            uuid.NewV1().String(),
			AccountID:     account.ID,
			TransactionID: uuid.NewV4().String(),
			JobNumber:     account.EmployeeID,
			CardNumber:    account.Customer.CardID,
			CustomerID:    account.CustomerID,
			CustomerName:  account.Customer.Name,
			Type:          lib.TypeSubsidy,
			Amounts:       -account.Virtual,
			Status:        lib.StatusRefund,
			Operator:      lib.OperatorSys,
			RealMoney:     account.Virtual,
			Balance:       account.Balance,
			DealTime:      dealTime,
			Remark:        "开卡补助",
		}
		rc.ShowDate = timeNow[0:10]
		rc.UpdateTime = timeNow
		if err := utils.WriteLedger(rc, stub, lib.RecordKey, []string{rc.ID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		account.VirtualTime = timeNow
	}
	account.Customer.CardNumber = "00" + account.EmployeeID
	// 账户信息写入账本
	if err := utils.WriteLedger(account, stub, lib.AccountsKey, []string{account.ID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	accounts, err := json.Marshal(account)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(accounts)
}

//查询房地产(可查询所有，也可根据所有人查询名下房产)
func QueryAccountsList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var realEstateList []lib.Accounts
	results, err := utils.GetStateByPartialCompositeKeys2(stub, lib.AccountsKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var account lib.Accounts
			err := json.Unmarshal(v, &account)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryRealEstateList-反序列化出错: %s", err))
			}
			realEstateList = append(realEstateList, account)
		}
	}
	accountsByte, err := json.Marshal(realEstateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRealEstateList-序列化出错: %s", err))
	}
	return shim.Success(accountsByte)
}
