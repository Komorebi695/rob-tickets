package mysql

import (
	"gorm.io/gorm"
	"robTickets/src/models"
)

const (
	orderPaid int8 = 1
)

// DescTicketNumber ,减少票数量
func DescTicketNumber(ticketID string, quantity int) error {
	if err := db.Exec("update `tickets` set `ticket_count`=`ticket_count` - ? where `ticket_id` = ? and ticket_count >= ?;", quantity, ticketID, quantity).Error; nil != err {
		return err
	}
	return nil
}

// IncrTicketNumber ,增加票数量
func IncrTicketNumber(ticketID string, count int) error {
	if err := db.Exec("update `tickets` set `ticket_count`=`ticket_count` + ? where `ticket_id` = ?;", count, ticketID, count).Error; err != nil {
		return err
	}
	return nil
}

// CreateOrder 生成订单
func CreateOrder(order models.TicketOrder) error {
	if err := db.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

// CancelOrder 取消/更新 订单
func CancelOrder(orderID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		order := models.TicketOrder{}
		if err := tx.Raw("select `status`,`ticket_id`,`quantity` from `ticket_orders` where `order_id` = ? for update;", orderID).Scan(&order).Error; err != nil {
			return err
		}

		// 已经支付
		if order.Status == orderPaid {
			// 需要退款
		}

		// 软删除订单
		if err := db.Exec("update `ticket_orders` set `deleted` = 1 where `order_id` = ?;", orderID).Error; err != nil {
			return err
		}

		if err := tx.Exec("update `tickets` set `ticket_count`=`ticket_count` + ? where `ticket_id` = ?;", order.Quantity, order.TicketId).Error; err != nil {
			return err
		}
		return nil
	})
}
