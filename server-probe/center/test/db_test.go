package test

import (
	"config-manager/center/db"
	"config-manager/center/global"
	"config-manager/common/model/node"
	"config-manager/common/utils"
	"fmt"
	"testing"
	"time"
)

func TestDB(t *testing.T) {

	db.InitDB()

	probeId := "Intel_Xeon"

	n := node.Node{}

	global.DB.Create(&node.Node{
		Id:      utils.UUID(),
		Name:    probeId,
		Address: "123",
		RegTime: time.Now(),
		Online:  true,
	})

	global.DB.Debug().Find(&n, probeId)

	fmt.Println(n)
	global.DB.Debug().Model(&n).Where("id = ?", probeId).Update("online", true)
	global.DB.Debug().Find(&n, probeId)

	fmt.Println(n)

}
