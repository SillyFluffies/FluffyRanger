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
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"

	fluf "fluffyranger/internal"
	flufCommands "fluffyranger/internal/commands"
	flufComponents "fluffyranger/internal/components"
	flufEvents "fluffyranger/internal/events"
	"fluffyranger/pkg/logger"
)

func main() {
	cfg, err := fluf.LoadConfig("config.toml")
	if err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
		os.Exit(-1)
	}

	logger.SetupLogger(cfg.Log.Format, &slog.HandlerOptions{
		Level:     cfg.Log.Level,
		AddSource: cfg.Log.AddSource,
	})
	slog.Info("Starting bot-template...")

	b, err := disgo.New(os.Getenv("token"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentMessageContent)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
		flufCommands.Setup(),
		flufEvents.Setup(),
		flufComponents.Setup(),
	)
	if err != nil {
		slog.Error("Failed to create bot", slog.Any("err", err))
		os.Exit(-1)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		b.Close(ctx)
	}()

	slog.Info("Syncing commands", slog.Any("guild_ids", cfg.Bot.DevGuilds))
	if err = handler.SyncCommands(b, flufCommands.Cmds, cfg.Bot.DevGuilds); err != nil {
		slog.Error("Failed to sync commands", slog.Any("err", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = b.OpenGateway(ctx); err != nil {
		slog.Error("Failed to open gateway", slog.Any("err", err))
		os.Exit(-1)
	}

	slog.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	slog.Info("Shutting down bot...")
}

func OnReady(e *events.Ready) {
	slog.Info("bot-template ready")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Client().SetPresence(ctx, gateway.WithListeningActivity("you"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
		slog.Error("Failed to set presence", slog.Any("err", err))
	}
}
