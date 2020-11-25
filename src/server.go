package ttrn

// Run runs the server application.
func Run(cfg Config, doDebug bool) {
	web := NewWeb(cfg, doDebug)
	web.Run()
}
