package logic

func init() {
	HandlerMap = make(map[string]Handler)
	HandlerMap["FluctuateHandler"] = &FluctuateHandler{}
}
