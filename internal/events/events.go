package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func Setup() bot.ConfigOpt {
	return bot.WithEventListeners(MessageHandler(), OnReady())
}

func MessageHandler() bot.EventListener {
	return bot.NewListenerFunc(func(e *events.MessageCreate) {
		// TODO: handle message
	})
}

func OnReady() bot.EventListener {
	return bot.NewListenerFunc(func(e *events.Ready) {
		slog.Info("bot-template ready")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Client().SetPresence(ctx, gateway.WithListeningActivity("you"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
			slog.Error("Failed to set presence", slog.Any("err", err))
		}
	})
}
