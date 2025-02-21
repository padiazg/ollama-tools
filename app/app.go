package app

import "github.com/padiazg/ollama-tools/models/settings"

type Application struct {
	Settings *settings.Settings
}

func (a *Application) Use(app MicroAppInterface) {
	app.ConfigureApplication(a)
}
