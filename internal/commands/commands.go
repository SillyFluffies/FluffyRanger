package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
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

func Sync(b *bot.Client, devGuilds []snowflake.ID) {
	slog.Info("Syncing commands", slog.Any("guild_ids", devGuilds))
	
	if err := handler.SyncCommands(b, Cmds, devGuilds); err != nil {
		slog.Error("Failed to sync commands", slog.Any("err", err))
	}
}
