package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"

	"github.com/sillyfluffies/fluffyranger/internal/commands"
	"github.com/sillyfluffies/fluffyranger/internal/components"
	"github.com/sillyfluffies/fluffyranger/internal/config"
	"github.com/sillyfluffies/fluffyranger/internal/events"
)

var (
	loggerFormat = "custom"
	loggerOpts   = slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	
	devGuilds = []snowflake.ID{}
	token     = os.Getenv("token")
	
	intents = gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentMessageContent
	caches = cache.FlagGuilds | cache.FlagMembers 
)

func main() {
	config.SetupLogger(loggerFormat, &loggerOpts)
	slog.Info("Starting bot")

	b, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(intents)),
		bot.WithCacheConfigOpts(cache.WithCaches(caches)),
		commands.Setup(),
		events.Setup(),
		components.Setup(),
	)
	if err != nil {
		slog.Error("Failed to create bot", slog.Any("err", err))
		os.Exit(-1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = b.OpenGateway(ctx); err != nil {
		slog.Error("Failed to open gateway", slog.Any("err", err))
		b.Close(ctx)
		os.Exit(-1)
	}

	commands.Sync(b, devGuilds)
	slog.Info("Bot is running. Press CTRL-C to exit")
	
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	slog.Info("Shutting down bot...")
}
