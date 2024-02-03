package test

import (
	"config-manager/center/action"
	"config-manager/center/db"
	"config-manager/center/global"
	"config-manager/center/model"
	"config-manager/center/service"
	"config-manager/center/service/chain"
	"fmt"
	"testing"
	"time"
)

func TestDispatch(t *testing.T) {

	prepareData()
	chainId := "38baba99-6051-494f-bb14-e40287713ac5"

	var nodes []model.Node
	global.DB.Where("chain_id = ?", chainId).Find(&nodes)

	var edges []model.Edge
	global.DB.Where("chain_id = ?", chainId).Find(&edges)

	var nodeIds []string

	for _, node := range nodes {
		nodeIds = append(nodeIds, node.Id)
	}

	var bindings []model.ShortcutNodeBinding

	global.DB.Where("node_id in ?", nodeIds).Find(&bindings)

	var shortcutIds []string

	for _, binding := range bindings {
		shortcutIds = append(shortcutIds, binding.ShortcutId)
	}

	var shortcuts []model.Shortcut

	global.DB.Where("id in ?", shortcutIds).Find(&shortcuts)

	params := action.ChainExecParams{
		Shortcuts: shortcuts,
		Nodes:     nodes,
		Edges:     edges,
		Bindings:  bindings,
	}

	dispatcher, err := action.NewDispatcher(params)

	if err != nil {
		t.Fatal(err)
	}

	for {
		err := dispatcher.Next()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(dispatcher.Out)
		if dispatcher.Ok {
			break
		}
		time.Sleep(time.Second)
	}

}

func prepareData() {

	db.InitDB()
	service.Init()

	chainId := "38baba99-6051-494f-bb14-e40287713ac5"

	service.ChainCRUD.SaveChain(&model.Chain{
		Id:          chainId,
		Name:        "test",
		Description: "test",
	})

	service.ChainCRUD.SaveNode(&model.Node{
		Id:          "1",
		Name:        "1",
		ChainId:     chainId,
		Description: "1",
		Root:        true,
	})

	service.ChainCRUD.SaveNode(&model.Node{
		Id:          "2",
		Name:        "2",
		ChainId:     chainId,
		Description: "2",
		Root:        false,
	})

	service.ShortcutCRUD.SaveShortcut(&model.Shortcut{
		Id:          "1",
		Name:        "1",
		Description: "1",
		Type:        model.ONE_LINE,
		CreateTime:  time.Now(),
		Timeout:     1000,
		JustRun:     false,
		Payload:     "echo 1",
		ProbeId:     "mac-os",
	})

	service.ShortcutCRUD.SaveShortcut(&model.Shortcut{
		Id:          "2",
		Name:        "2",
		Description: "2",
		Type:        model.ONE_LINE,
		CreateTime:  time.Now(),
		Timeout:     1000,
		JustRun:     false,
		Payload:     "echo 2",
		ProbeId:     "2",
	})

	service.ChainCRUD.LinkNode(chain.ConnectTwoNodesParams{
		SourceId: "1",
		TargetId: "2",
		ChainId:  chainId,
		Type:     model.SUCCESS,
	})

	service.ChainCRUD.BindShortcut("1", "1")
	service.ChainCRUD.BindShortcut("2", "2")

	chainInfo, _ := service.ChainCRUD.ChainInfo(chainId)

	fmt.Println(chainInfo.Chain)
	fmt.Println(chainInfo.Nodes)
	fmt.Println(chainInfo.Edges)

}
