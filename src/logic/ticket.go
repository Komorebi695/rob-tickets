package logic

import (
	"encoding/json"
	"gopkg.in/Shopify/sarama.v1"
	"robTickets/src/dao/redis"
	"robTickets/src/models"
	"robTickets/src/util"
	"time"
)

const (
	waitTime = time.Minute * 10
)

// BuyTicket ,购票 (版本1.0)
//func BuyTicket(userID string, ticketID string, quantity int) (models.TicketOrder, bool, error) {
//	order := models.TicketOrder{
//		OrderId:  util.GenSnowID(),
//		TicketId: ticketID,
//		UserId:   userID,
//		Quantity: quantity,
//	}
//
//	// MySQL票减
//	if err := mysql.DescTicketNumber(ticketID, quantity); err != nil {
//		return models.TicketOrder{}, false, err
//	}
//	// 增加订单
//	if err := mysql.CreateOrder(order); err != nil {
//		return models.TicketOrder{}, false, err
//	}
//
//	// 开启定时任务，若5分钟内未支付，则自动取消订单
//	go func() {
//		timer := time.NewTimer(waitTime)
//		// 阻塞5分钟
//		_ = <-timer.C
//		if err := mysql.CancelOrder(order.OrderId); err != nil {
//			log.Fatal(err)
//		}
//	}()
//	return order, true, nil
//}

// BuyTicket ,购票
func BuyTicket(userID string, ticketID string, quantity int) (models.TicketOrder, bool, error) {
	//// 校验ticketID是否存在;
	//exist, err := redis.Operation().ExistKey(ticketID)
	//if err != nil || !exist {
	//	return models.TicketOrder{}, false, err
	//}

	// redis缓存中票数减少，返回当前票数
	curNumber, err := redis.Operation().DecrTicket(ticketID, int64(quantity))
	if err != nil {
		return models.TicketOrder{}, false, err
	}

	//fmt.Println(endTime.Sub(beginTime))
	// 票数不够，购买失败。
	if curNumber < 0 {
		return models.TicketOrder{}, false, nil
	}

	order := models.TicketOrder{
		OrderId:  util.GenSnowID(),
		TicketId: ticketID,
		UserId:   userID,
		Quantity: quantity,
	}

	mTicket := models.MTicket{
		TicketID: ticketID,
		Quantity: quantity,
		Type:     decr,
	}

	cancelOrder := models.MCancelOrder{
		OrderID: order.OrderId,
	}

	// 序列化
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return models.TicketOrder{}, false, err
	}
	mTicketBytes, err := json.Marshal(mTicket)
	if err != nil {
		return models.TicketOrder{}, false, err
	}
	cancelOrderBytes, err := json.Marshal(cancelOrder)
	if err != nil {
		return models.TicketOrder{}, false, err
	}
	// 插入订单
	SendMessage(updateTicketNumberTopic, sarama.ByteEncoder(mTicketBytes))
	// 更新票数
	SendMessage(insertOrderTopic, sarama.ByteEncoder(orderBytes))

	// 开启定时任务，若10分钟内未支付，则自动取消订单
	go func(cancelOrder []byte) {
		timer := time.NewTimer(waitTime)
		// 阻塞10分钟
		_ = <-timer.C
		SendMessage(cancelOrderTopic, sarama.ByteEncoder(cancelOrder))
	}(cancelOrderBytes)
	return order, true, nil
}

// CancelOrder ,取消订单（主动）
func CancelOrder(orderID string) error {
	// todo 用户主动取消订单
	return nil
}

// SetRedisTicketNumber ,设置redis中的初始票数量
func SetRedisTicketNumber(ticketID string, number int) error {
	return redis.Operation().SetTicketNumber(ticketID, number)
}
