package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var ctx = context.Background()

const TicketNumberPer = "ticketNumber"

type redisFunc interface {
	SetTicketNumber(ticketID string, ticketNumber int) error
	IncrTicket(ticketID string, incr int64) (int64, error)
	DecrTicket(ticketID string, decr int64) (int64, error)
	ExistKey(ticketID string) (bool, error)
	GetTicketValue(ticketID string) (int, error)
}

type redisBase struct{}

var _ redisFunc = (*redisBase)(nil)

func Operation() *redisBase {
	return &redisBase{}
}

func (c *redisBase) SetTicketNumber(ticketID string, ticketNumber int) error {
	if err := rdb.Set(ctx, fmt.Sprintf("%s:%s", TicketNumberPer, ticketID), ticketNumber, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (c *redisBase) IncrTicket(ticketID string, incr int64) (int64, error) {
	number, err := rdb.IncrBy(ctx, fmt.Sprintf("%s:%s", TicketNumberPer, ticketID), incr).Result()
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (c *redisBase) DecrTicket(ticketID string, decr int64) (int64, error) {
	var decrBy = redis.NewScript(`
	local key = KEYS[1]
	local change = ARGV[1]
	
    local value = redis.call("GET",key)
	if value < change then
		return -1
	end

	value = value - change
	redis.call("SET",key,value)
	return value
	`)
	keys := []string{fmt.Sprintf("%s:%s", TicketNumberPer, ticketID)}
	values := []interface{}{decr}
	num, err := decrBy.Run(ctx, rdb, keys, values).Int64()
	if err != nil {
		return 0, err
	}
	//number, err := rdb.DecrBy(ctx, fmt.Sprintf("%s:%s", TicketNumberPer, ticketID), decr).Result()
	//if err != nil {
	//	return 0, err
	//}
	return num, nil
}

func (c *redisBase) ExistKey(ticketID string) (bool, error) {
	exist, err := rdb.Exists(ctx, fmt.Sprintf("%s:%s", TicketNumberPer, ticketID)).Result()
	if err != nil {
		return false, err
	}
	return exist == int64(1), nil
}

func (c *redisBase) GetTicketValue(ticketID string) (int, error) {
	curStr, err := rdb.Get(ctx, fmt.Sprintf("%s:%s", TicketNumberPer, ticketID)).Result()
	cur, _ := strconv.Atoi(curStr)
	return cur, err
}
