package snowflake

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

// Generator holds a snowflake node for generating IDs.
type Generator struct {
	*snowflake.Node // The snowflake node to use for generating IDs.
}

// NewGenerator creates a new Generator instance with the specified nodeID.
func NewGenerator(nodeID int64) *Generator {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("unable to generate snowflake: %v", err)
	}
	return &Generator{
		Node: node,
	}
}
