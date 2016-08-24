package main

import (
	"fmt"
	"time"

	"github.com/tinycedar/lily/common"
	"github.com/tinycedar/lily/core"
)

var batcher *core.Batcher

func main() {
	initBatcher()
	fmt.Println("batcher in main: ", batcher)
}

func initBatcher() {
	var err error
	batcher, err = core.New(time.Millisecond * 300)
	if err != nil {
		common.Error("Fail to initialize batcher")
	}
}
