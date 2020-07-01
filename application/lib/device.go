package lib

type GetTimesResponse struct {
	Status string `json:"Status"`
	Msg    string `json:"Msg"`
	Time   string `json:"Time"`
}

type ConsumTransactionsCond struct {
	MachineCode string        `json:"machine_code"`
	Order       string        `json:"Order"`
	CardNo      string        `json:"CardNo"`
	Mode        string        `json:"Mode"`
	PayType     string        `json:"PayType"`
	Amount      string        `json:"Amount"`
	Menus       []interface{} `json:"Menus"`
	Code        string        `json:"Code"`
	QR          string        `json:"QR"`
	QROrder     string        `json:"QROrder"`
	PassWord    string        `json:"PassWord"`
	Name        string        `json:"Name"`
	Money       string        `json:"Money"`
	Subsidy     string        `json:"Subsidy"`
	Finger      string        `json:"Finger"`
}

type ConsumTransactionsResponse struct {
	StartTime    string `json:"StartTime"`
	EndTime      string `json:"EndTime"`
	Status       string `json:"Status"`
	Msg          string `json:"Msg"`
	Name         string `json:"Name"`
	Money        string `json:"Money"`
	Subsidy      string `json:"Subsidy"`
	Times        string `json:"Times"`
	Integral     string `json:"Integral"`
	Dept         string `json:"Dept"`
	Discount     string `json:"Discount"`
	QR           string `json:"QR"`
	CardType     string `json:"CardType"`
	InTime       string `json:"InTime"`
	OutTime      string `json:"OutTime"`
	Amount       string `json:"Amount"`
	QROrder      string `json:"QROrder"`
	PayResults   string `json:"PayResults"`
	Level        string `json:"Level"`
	CardNumber   string `json:"CardNumber"`
	QRCodeNumber string `json:"QRCodeNumber"`
	VoiceID      string `json:"VoiceID"`
	Text         string `json:"Text"`
}
