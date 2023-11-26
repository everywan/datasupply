package constant

var Closedchan = make(chan struct{})

func init() {
	close(Closedchan)
}
