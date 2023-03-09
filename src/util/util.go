package util

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"strings"
)

const (
	nodeID int64 = 695 // 大于0且小于int64，否则报错
)

// SerialNumber 获取uuid后12位
func SerialNumber() string {
	id := uuid.New()
	// da90336c-7ea8-4b00-8fb9-6cf90e5cec5d  ->  6cf90e5cec5d
	return strings.Split(id.String(), "-")[4]
}

// GenSnowID 生成ID时会上锁，确保不重复 11位
func GenSnowID() string {
	node, _ := snowflake.NewNode(nodeID)
	return node.Generate().Base58()
}
