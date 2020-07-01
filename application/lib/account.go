package lib

type Account struct {
	ID             string   `json:"id" gorm:"type:varchar(64);primary_key"`
	Customer       Customer `json:"owner" gorm:"ForeignKey:customer_id"`
	CustomerID     string   `json:"customer_id" gorm:"type:varchar(64);index"`                  //关联的卡id
	Balance        int      `json:"balance" gorm:"type:int(11)"`                                //余额
	Status         string   `json:"status" gorm:"type:char(1)"`                                 //账户状态 1正常 2冻结
	Name           string   `json:"name" gorm:"type:varchar(15)"`                               //姓名
	EmployeeID     string   `json:"employee_id" gorm:"type:varchar(11)"`                        //人员编号
	Department     string   `json:"department" gorm:"type:varchar(64)"`                         //部门ID
	CardFee        int      `json:"card_fee" gorm:"column:card_fee;type:int(11)"`               //卡费
	EmployeeStatus int      `json:"employee_status" gorm:"column:employee_status;type:int(11)"` //员工在职状态（1在职 2离职 ）
	OldNO          string   `json:"old_no" gorm:"type:varchar(16)"`
	StaffGroupID   string   `json:"staff_group_id" gorm:"type:varchar(64)"`
	Cash           int      `json:"cash" gorm:"type:int(11)"`
	Virtual        int      `json:"virtual" gorm:"type:int(11)"`
	VirtualTime    string   `json:"virtual_time" gorm:"type:varchar(32)"`
	CreateTime     string   `json:"create_time" gorm:"type:varchar(32)"`
	UpdateTime     string   `json:"update_time" gorm:"type:varchar(32)"`
}

type UpdateAccountsRequest struct {
	ID             string `json:"id"`
	EmployeeID     string `json:"employee_id"`      //员工编号
	CardID         string `json:"card_id"`          //卡号
	StaffGroupName string `json:"staff_group_name"` //人员组名称
	Balance        int    `json:"balance"`          //余额
	Sums           int    `json:"sums"`             //上次补助金额
	SumsTime       string `json:"sums_time"`        //上次补助时间
	Department     string `json:"department"`       //部门
	Status         string `json:"status"`           //账户状态
	Name           string `json:"name"`             //名称
	IDCard         string `json:"id_card"`          //身份证号码
	Cellphone      string `json:"cellphone"`        //电话
	CreateTime     string `json:"create_time"`      //创建时间
	CardFee        int    `json:"card_fee"`         //卡费
	EmployeeStatus int    `json:"employee_status"`  //员工状态
	DepartmentCode string `json:"department_code"`
	CardStatus     string `json:"card_status"` //卡状态
}

type Tem struct {
	ID      string `json:"id" `
	Cash    int    `json:"cash" `
	Virtual int    `json:"virtual"`
}
