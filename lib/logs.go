package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type LogType string

const (
	LogTypeError = "error"
	LogTypeInfo  = "info"
)

const (
	LogErrorColor = 15548997
	LogInfoColor  = 3447003
)

type LogProps struct {
	Type        LogType
	Title       string
	Description string
	Message     string
}

func SendLog(sendlog *LogProps) {
	data := map[string]any{}

	var color int = 0
	if sendlog.Type == LogTypeError {
		color = LogErrorColor
	}
	if sendlog.Type == LogTypeInfo {
		color = LogInfoColor
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Cordy Logs",
		},
		Title:       sendlog.Title,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: sendlog.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Details",
				Value: sendlog.Message,
			},
			{
				Name:  "Timestamp",
				Value: fmt.Sprintf("<t:%d:R>", time.Now().Unix()),
			},
		},
		Color: color,
	}

	data["embeds"] = []*discordgo.MessageEmbed{embed}

	var resp map[string]interface{}
	if err := request.Post(WEBHOOKLOGS_API, data, &resp); err != nil {
		log.Printf("Error sending logs to webhook: %v\n", err)
	}
}
