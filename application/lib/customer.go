package lib

type Customer struct {
	ID             string         `json:"id" gorm:"type:varchar(64);primary_key"`
	Name           string         `json:"name" gorm:"type:varchar(15)"`
	IDCard         string         `json:"idcard" gorm:"type:char(18)"`        //身份证
	EmployeeID     string         `json:"employeeID" gorm:"type:varchar(11)"` //工号
	Cellchone      string         `json:"cellphone" gorm:"type:varchar(11)"`  //手机号
	AccountID      string         `json:"accountID" gorm:"type:varchar(64);index"`
	CardID         string         `json:"cardID" gorm:"type:varchar(64)"`               //卡号
	CardFee        int            `json:"card_fee" gorm:"column:card_fee;type:int(11)"` //卡费
	CardNumber     string         `json:"cardNumber" gorm:"type:varchar(32)"`           //卡内码
	Department     string         `json:"department" gorm:"type:varchar(64)"`           //部门ID
	CardCreateTime string         `json:"cardcreatetime" gorm:"type:varchar(64)"`       //卡创建时间
	CardEndTime    string         `json:"cardendtime" gorm:"type:varchar(64)"`          //卡过期时间
	Cardstatus     string         `json:"cardstatus"`                                   //卡状态 1正常 2挂失
	Types          int            `json:"types" gorm:"type:int(1);column:types"`        //卡类型（1新办卡  2补卡）
}
