package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/global"
	"apps/center/model"
	"apps/center/service"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ShortcutCreate 创建快捷方式
// Params: probeId 探针id，name 快捷方式名称，timeout 超时时间，justRun 是否只运行一次，payload 脚本内容，type 脚本类型，description 描述
// Response: 200，“success”， true
// Err: 200, “fail”， err
func ShortcutCreate(c *gin.Context) {
	p := param.CreateShortcutParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	ok := checkProbeId(p.ProbeId)

	if !ok {
		c.JSON(http.StatusOK, response.StringFail("can not find a probe with id = "+p.ProbeId))
		return
	}

	sc := model.Shortcut{}
	sc.ProbeId = p.ProbeId
	sc.Name = p.Name
	sc.Timeout = p.Timeout
	sc.JustRun = p.JustRun
	sc.CreateTime = time.Now()
	sc.Payload = p.Payload
	sc.Type = p.Type
	sc.Description = p.Description
	sc.Args = p.Args

	err = service.ShortcutCRUD.SaveShortcut(&sc)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[bool](true))
}

// ListShortcuts 列出指定的所有快捷方式
// Params: probeId 探针id
// Response: 200，“success”， []model.Shortcut
// Err: 200, “fail”， err
func ListShortcuts(c *gin.Context) {

	probeId, b := c.GetQuery("probeId")

	if !b {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	ok := checkProbeId(probeId)

	if !ok {
		c.JSON(http.StatusOK, response.StringFail("can not find a probe with id = "+probeId))
		return
	}

	shortcuts, err := service.ShortcutCRUD.ListShortcuts(probeId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(shortcuts))

}

// ShortcutGroup 快捷方式分组汇总
// Params: 无
// Response: 200，“success”， map[string][]model.Shortcut(key:探针id，value:快捷方式数组)
func ShortcutGroup(c *gin.Context) {

	groups, err := service.ShortcutCRUD.ShortcutGroup()

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(groups))

}

// RunShortcut 运行快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， {ok: true, out: out}
// Err: 200, “fail”， {ok: false, out: err}
func RunShortcut(c *gin.Context) {

	var runScriptParams param.RunScriptParams

	err := c.ShouldBindJSON(&runScriptParams)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	stdout, stderr, state, err := service.ShortcutRUN.Run(runScriptParams.ShortcutId)

	if err != nil {
		c.JSON(http.StatusOK, response.Success(gin.H{
			"ok":  false,
			"out": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"state":  state,
		"stdout": stdout,
		"stderr": stderr,
	}))
}

// AsyncRunShortcut realtime run shortcut
// push result using websocket
// TODO test
func AsyncRunShortcut(c *gin.Context) {

	scriptId, ok := c.GetQuery("shortcutId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	w := c.Writer
	flusher, ok := w.(http.Flusher)

	if !ok {
		c.JSON(http.StatusOK, response.StringFail("unknown error"))
		return
	}

	streams, err := service.ShortcutRUN.AsyncRun(scriptId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	for stream := range streams {

		log.Println(stream.Content)

		b, err := json.Marshal(stream)

		if err != nil {
			continue
		}

		w.Write(b)
		w.Write([]byte("\n"))
		flusher.Flush()
	}
}

// ShortcutRunHistory 快捷方式运行历史
// Params: shortcutId 快捷方式id
// Response: 200，“success”， []model.ShortcutExecLog
// Err: 200, “fail”， err
func ShortcutRunHistory(c *gin.Context) {

	shortcutId, ok := c.GetQuery("shortcutId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	logs, err := service.ShortcutRUN.RunHistory(shortcutId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(logs))
}

// DeleteShortcut 删除快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， true
// Err: 200, “fail”， err
func DeleteShortcut(c *gin.Context) {

	id, ok := c.GetQuery("shortcutId")

	if !(ok) {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}
	sc := model.Shortcut{}
	service.ShortcutCRUD.Db.Model(&model.Shortcut{}).Where("id = ?", id).First(&sc)

	if sc.Type == model.SCRIPT {
		// 删除脚本文件
		service.FS_OPERATION.DeleteFile(sc.ProbeId, sc.Payload)
	}

	err := service.ShortcutCRUD.RemoveShortcut(id)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// UpdateShortcut 更新快捷方式
// Params: shortcutId 快捷方式id
// Response: 200，“success”， true
// Err: 200, “fail”， err
func UpdateShortcut(c *gin.Context) {
	var shortcut = model.Shortcut{}

	err := c.ShouldBindJSON(&shortcut)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ShortcutCRUD.UpdateShortcut(&shortcut)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())

}

func checkProbeId(probeId string) bool {
	//检测probe是否存在
	var count int64
	global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probeId).Count(&count)
	if count < 1 {
		return false
	}

	return true
}
