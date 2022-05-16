package util

type CallbackQueue chan func()

func (queue *CallbackQueue) NextFrame(cb func()) {
	(*queue) <- func() {
		(*queue) <- cb
	}
}

func (queue *CallbackQueue) Update() {
	select {
	case cb := <-(*queue):
		cb()
	default:
		break
	}
}
