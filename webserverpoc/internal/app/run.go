package app

import gowebserver "github.com/edwardmccormick/gowebserver"

func Run() error {
	config, err := gowebserver.LoadRuntimeConfig()
	if err != nil {
		return err
	}

	app, err := gowebserver.NewApp(config)
	if err != nil {
		return err
	}

	return app.Router().Run(gowebserver.DefaultBindAddress())
}
