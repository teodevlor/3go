package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-structure/config"
	"go-structure/pkg/telegram"
)

func main() {
	if err := config.Load("./config", "dev.yml"); err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if config.Cfg == nil {
		log.Fatal("config is nil after Load")
	}

	tgCfg := config.Cfg.TelegramConfig
	if tgCfg.BotToken == "" || tgCfg.ChatID == 0 {
		log.Fatal("telegram config is empty (bot_token/chat_id)")
	}

	client := telegram.NewClient(tgCfg.BotToken, tgCfg.ChatID)

	ctx := context.Background()
	msg := "🔥 Telegram test message from go-structure"

	if len(os.Args) > 1 {
		msg = os.Args[1]
	}

	if err := client.SendMessage(ctx, msg); err != nil {
		log.Fatalf("send telegram message failed: %v", err)
	}

	fmt.Println("Telegram test message sent successfully.")
}
