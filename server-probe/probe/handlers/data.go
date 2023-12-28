package handlers

import (
	"config-manager/common/message"
	. "config-manager/common/message/data"
	//shortcutservice "config-manager/common/shortcut"
	fileservice "config-manager/probe/service/fs"
	"encoding/json"
)

var DataHandlePolicy = map[message.DataType]func(data []byte) ([]byte, error){
	message.READDIR:    handleReadDir,
	message.READFILE:   handleReadFile,
	message.MODIFYFILE: handleModifyFile,
	message.CREATEFILE: handleCreateFile,
	//message.CREATE_ONELINE_SHORTCUT: handleCreateOneLineShortcut,
	//message.CREATE_SCRIPT_SHORTCUT:  handleCreateScriptShortcut,
	//message.DELETE_SHORTCUT:         handleDeleteShortcut,
}

func readData[T any](data []byte) (T, error) {

	res := new(T)

	err := json.Unmarshal(data, res)
	if err != nil {
		return *res, err
	}

	return *res, nil
}

func handleReadDir(data []byte) ([]byte, error) {
	//将msg.Data解析为对应类型
	dr, err := readData[DirRead](data)

	if err != nil {
		return nil, err
	}

	//调用对应处理方法
	dirReadResponse, err := fileservice.DirRead(dr)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(dirReadResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func handleReadFile(data []byte) ([]byte, error) {
	fileRead, err := readData[FileRead](data)

	if err != nil {

		return nil, err
	}

	fileReadResponse, err := fileservice.FileRead(fileRead)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(fileReadResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func handleModifyFile(data []byte) ([]byte, error) {

	mf, err := readData[FileModify](data)

	if err != nil {
		return nil, err
	}

	result, err := fileservice.FileModify(mf)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func handleCreateFile(data []byte) ([]byte, error) {

	fc, err := readData[FileCreate](data)

	if err != nil {
		return nil, err
	}

	fileCreateResponse, err := fileservice.FileCreate(fc)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(fileCreateResponse)

	if err != nil {
		return nil, err
	}

	return bytes, nil

}

//
//func handleCreateOneLineShortcut(data []byte) ([]byte, error) {
//
//	sc, err := readData[OneLineShortcutCreate](data)
//
//	if err != nil {
//		return nil, err
//	}
//
//	shortcut := shortcutservice.NewOneLineShortcut(sc.Name, sc.Command, sc.Timeout, sc.JustRun, sc.Description)
//
//	db.DB.Create(&shortcut)
//
//	resp := OneLineShortcutCreateResp{
//		Ok:         true,
//		ShortcutId: shortcut.Id,
//		CreateTime: shortcut.CreateTime,
//	}
//
//	bytes, err := json.Marshal(resp)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return bytes, nil
//}
//
//func handleCreateScriptShortcut(data []byte) ([]byte, error) {
//
//	sc, err := readData[ScriptShortcutCreate](data)
//
//	if err != nil {
//		return nil, err
//	}
//
//	scriptShortcut, err := shortcutservice.NewScriptShortcut(sc.Name, sc.Path, sc.Timeout, sc.JustRun, sc.Description)
//
//	db.DB.Create(&scriptShortcut)
//
//	if err != nil {
//		return nil, err
//	}
//
//	resp := ScriptShortcutCreateResp{
//		ShortcutId: scriptShortcut.Id,
//		Ok:         true,
//		CreateTime: scriptShortcut.CreateTime,
//	}
//
//	bytes, err := json.Marshal(resp)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return bytes, nil
//}
//
//func handleDeleteShortcut(data []byte) ([]byte, error) {
//
//	sd, err := readData[ShortcutDelete](data)
//
//	if err != nil {
//		return nil, err
//	}
//
//	resp := ShortcutDeleteResp{}
//
//	if sd.Type == ONE_LINE_SHORTCUT {
//		db.DB.Delete(&shortcutservice.OneLineShortcut{},sd.Id)
//		resp.OK = true
//	} else if sd.Type == SCRIPT_SHORTCUT {
//		db.DB.Delete(&shortcutservice.ScriptShortcut{},sd.Id)
//		resp.OK = true
//	} else {
//		resp.OK = false
//	}
//
//	bytes, err := json.Marshal(resp)
//
//	return bytes, nil
//}
