package snowflake

import (
	"log"
	"sync/atomic"

	"github.com/bwmarrin/snowflake"
)

var globalIDGenerator atomic.Pointer[snowflake.Node]

func SetGlobalIDGenerator(nodeNum int64) {
	node, err := snowflake.NewNode(nodeNum)
	if err != nil {
		log.Fatalf("unable to generate snowflake: %v", err)
	}
	globalIDGenerator.Store(node)
}

func GenerateID() int64 {
	return globalIDGenerator.Load().Generate().Int64()
}
