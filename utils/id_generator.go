package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
)

func UUID(resource string) string {
	return uuid.NewSHA1(uuid.New(), strings2.ToByte(resource)).String()
}

func UniqueId(nodeNum int64) (int64, error) {
	node, err := snowflake.NewNode(nodeNum)

	if err != nil {
		logrus.Debugf("Unique Id generation failed: %v", err)
		return 0, err
	}

	id := node.Generate()

	return id.Int64(), nil
}

