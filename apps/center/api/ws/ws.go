package ws

var Ctx *Context

func Init() {
	Ctx = newContext()
	Ctx.setup()
}
