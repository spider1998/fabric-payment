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

	// 写入账本
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
