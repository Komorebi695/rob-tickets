package logic

import (
	"log"
	"robTickets/src/dao/mysql"
	"robTickets/src/models"
	"robTickets/src/util"
	"time"
)

const (
	waitTime = time.Minute * 10
)

// BuyTicket ,购票
func BuyTicket(userID string, ticketID string, quantity int) (models.TicketOrder, bool, error) {
	order := models.TicketOrder{
		OrderId:  util.GenSnowID(),
		TicketId: ticketID,
		UserId:   userID,
		Quantity: quantity,
	}

	// MySQL票减
	if err := mysql.DescTicketNumber(ticketID, quantity); err != nil {
		return models.TicketOrder{}, false, err
	}
	// 增加订单
	if err := mysql.CreateOrder(order); err != nil {
		return models.TicketOrder{}, false, err
	}

	// 开启定时任务，若5分钟内未支付，则自动取消订单
	go func() {
		timer := time.NewTimer(waitTime)
		// 阻塞5分钟
		_ = <-timer.C
		if err := mysql.CancelOrder(order.OrderId); err != nil {
			log.Fatal(err)
		}
	}()
	return order, true, nil
}

// CancelOrder ,取消订单（主动）
func CancelOrder(orderID string) error {
	// todo 用户主动取消订单
	return nil
}
