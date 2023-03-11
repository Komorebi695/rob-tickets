package models

type User struct {
	UserID   string `gorm:"type:string;size:64" json:"user_id"`  // 用户名
	Password string `gorm:"type:string;size:64" json:"password"` // 密码
}

type TicketOrder struct {
	ID       uint   `gorm:"primary_key"`
	OrderId  string `gorm:"type:string;size:64"`                                      // 订单编号
	UserId   string `gorm:"type:string;size:64"`                                      // 用户编号
	TicketId string `gorm:"type:string;size:64"`                                      // 票编号
	Quantity int    `gorm:"column：quantity"`                                          // 票数量
	Status   int8   `gorm:"column:status;type:tinyint;size:1;not null;default:(0)"`   // 订单状态
	Deleted  int8   `gorm:"column:deleted;;type:tinyint;size:1;not null;default:(0)"` //删除状态
}

type MTicket struct {
	TicketID string `json:"ticket_id"`
	Quantity int    `json:"quantity"`
	Type     int8   `json:"type"`
}

type MCancelOrder struct {
	OrderID string `json:"order_id"`
}
