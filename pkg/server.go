package pkg

// Run runs the server application.
func Run(cfg Config) {
	web := NewWeb(cfg)
	cfg.run()
}
