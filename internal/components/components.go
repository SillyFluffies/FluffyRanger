package components

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
)

func Setup() bot.ConfigOpt {
	h := handler.New()
	h.Component("/test-button/{id}", TestComponent)
	
	return bot.WithEventListeners(h)
}