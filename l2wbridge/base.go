package l2wbridge

//logic -> widget
type L2WHandler interface {
	WidgetEventHandler(string, string)
}

//widget -> logic
type W2LHandler interface {
	LogicEventHandler(string, string)
}
