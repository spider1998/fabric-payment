package lib

type Account struct {
	ID             string     `json:"id" gorm:"type:varchar(64);primary_key"`
	Customer       Customer   `json:"owner" gorm:"ForeignKey:customer_id"`
	CustomerID     string     `json:"-" gorm:"type:varchar(64);index"`                            //关联的卡id
	Balance        int        `json:"balance" gorm:"type:int(11)"`                                //余额
	Status         string     `json:"status" gorm:"type:char(1)"`                                 //账户状态 1正常 2冻结
	Name           string     `json:"-" gorm:"type:varchar(15)"`                                  //姓名
	EmployeeID     string     `json:"-" gorm:"type:varchar(11)"`                                  //人员编号
	Department     string     `json:"department" gorm:"type:varchar(64)"`                         //部门ID
	CardFee        int        `json:"card_fee" gorm:"column:card_fee;type:int(11)"`               //卡费
	EmployeeStatus int        `json:"employee_status" gorm:"column:employee_status;type:int(11)"` //员工在职状态（1在职 2离职 ）
	OldNO          string     `json:"old_no" gorm:"type:varchar(16)"`
	StaffGroupID   string     `json:"staff_group_id" gorm:"type:varchar(64)"`
	Cash           int        `json:"cash" gorm:"type:int(11)"`
	Virtual        int        `json:"virtual" gorm:"type:int(11)"`
	VirtualTime    string     `json:"virtual_time" gorm:"type:varchar(32)"`
	CreateTime     string     `json:"create_time" gorm:"type:varchar(32)"`
	UpdateTime     string     `json:"update_time" gorm:"type:varchar(32)"`
}


