package lib

const (
	RecordKey = "record-key"
)

// Record 记录
type Record struct {
	ID                string `json:"id" gorm:"column:id;type:varchar(36); primary_key"`
	AccountID         string `json:"accountID" gorm:"column:account_id;type:varchar(36)"`
	TransactionID     string `json:"transactionID" gorm:"column:transaction_id;type:varchar(36)"` // 交易流水ID
	MasterAccountName string `json:"masterAccountName" gorm:"column:master_account_name;type:varchar(100)"`
	Department        string `json:"department" gorm:"type:varchar(64)"` //部门
	JobNumber         string `json:"jobNumber" gorm:"column:job_number;type:varchar(64)"`
	CardNumber        string `json:"cardNumber" gorm:"column:card_number;type:varchar(64)"`
	CustomerID        string `json:"customerID" gorm:"column:customer_id;type:varchar(36)"`
	CustomerName      string `json:"customerName" gorm:"column:customer_name;type:varchar(100)"`
	Type              string `json:"type" gorm:"column:type;type:char(2)"`
	Amounts           int    `json:"amounts" gorm:"column:amounts;type:int" `
	Balance           int    `json:"balance" gorm:"column:balance;type:int" `
	MerchantID        string `json:"merchantID" gorm:"column:merchant_id;type:varchar(36)"`
	MerchantName      string `json:"merchantName" gorm:"column:merchant_name;type:varchar(100)"`
	DeviceID          string `json:"deviceID" gorm:"column:device_id;type:varchar(36)"`
	DeviceCode        string `json:"device_code"`
	DeviceName        string `json:"deviceName" gorm:"column:device_name;type:varchar(100)"`
	DealTime          int64  `json:"dealTime" gorm:"column:deal_time"`
	Remark            string `json:"remark" gorm:"column:remark;type:varchar(255)"`
	Status            string `json:"status" gorm:"column:status;type:char(2)"`
	Operator          string `json:"operator" gorm:"column:operator;type:varchar(64)"`
	VirtualMoney      int    `json:"virtualMoney" gorm:"column:virtual_money;type:int"`
	RealMoney         int    `json:"realMoney" gorm:"column:real_money;type:int"`
	VirtualBalance    int    `json:"virtualBalance" gorm:"column:virtual_balance;type:int"`
	RealBalance       int    `json:"realBalance" gorm:"column:real_balance;type:int"`
	ShowTime          string `json:"show_time" gorm:"column:show_time;type:varchar(20)"`
	ShowDate          string `json:"show_date" gorm:"column:show_date;type:varchar(10)"`
	UpdateTime        string `json:"update_time" gorm:"column:update_time;type:varchar(20)"`
}

const (
	// 卡消费
	TypeTransaction = "1"
	// 二维码消费
	typeQRCode = "2"
	//充值
	TypeRecharge = "3"
	//typeCancell 交易撤销
	typeCancell = "4"
	//typeRechargeRfund  充值退款
	typeRechargeRfund = "5"
	//卡费
	Typecard = "6"
	//补助
	TypeSubsidy = "7"
	//年底清零
	typeBalanceClear = "8"
	//纠错
	typeCorrection = "9"

	//typeBalanceDeduct  余额扣款
	// operatorSys 管理员
	OperatorSys = "sysadmin"
	// operatorPos pos机操作
	operatorPos = "pos"
	// statusRefund 可撤销交易
	StatusRefund = "1"
	// statusNoRefund 不可撤销交易
	statusNoRefund = "2"
	statusWaiting  = "3"
)
