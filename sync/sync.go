package sync

func Run(function func(), stop chan<- bool) {
	go function()
	stop <- true
}
