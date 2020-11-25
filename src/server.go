package ttrn

// Run runs the server application.
func Run(cfg Config) {
	web := NewWeb(cfg)
	web.Run()
}
