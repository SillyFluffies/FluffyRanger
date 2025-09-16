package events

import (
	"context"
	"fmt"
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
		fmt.Println("\033[38;5;220m" + `
  ______ _        __  __       _____
 |  ____| |      / _|/ _|     |  __ \
 | |__  | |_   _| |_| |_ _   _| |__) |__ _ _ __   __ _  ___ _ __
 |  __| | | | | |  _|  _| | | |  _  /  _' | '_ \ / _' |/ _ \ '__|
 | |    | | |_| | | | | | |_| | | \ \ (_| | | | | (_| |  __/ |
 |_|    |_|\__,_|_| |_|  \__, |_|  \_\__,_|_| |_|\__, |\___|_|
                          __/ |                   __/ |
                         |___/                   |___/
        `)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Client().SetPresence(ctx, gateway.WithListeningActivity("you"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
			slog.Error("Failed to set presence", slog.Any("err", err))
		}
	})
}
