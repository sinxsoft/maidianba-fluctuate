package logic

/**
所有实现Handler接口的类的map
*/

var HandlerMap map[string]Handler

type Handler interface {
	Handle() error
}
