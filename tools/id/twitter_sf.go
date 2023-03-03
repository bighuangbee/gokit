package id

import (
	"github.com/bwmarrin/snowflake"
)

type TwitterSF struct {
	node *snowflake.Node
}

func NewTwitterSF(workerId int64) (Generater, error) {
	node, err := snowflake.NewNode(workerId)
	if err != nil {
		return nil, err
	}
	return &TwitterSF{node: node}, nil
}

func (sf *TwitterSF) Generate() ID {
	return sf.node.Generate()
}
