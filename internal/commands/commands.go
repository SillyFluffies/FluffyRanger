package commands

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Cmds = []discord.ApplicationCommandCreate{
	test,
}

func Setup() bot.ConfigOpt {
	h := handler.New()
	h.Command("/test", TestHandler)
	h.Autocomplete("/test", TestAutocompleteHandler)
	
	return bot.WithEventListeners(h)
}
