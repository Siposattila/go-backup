package alert

import (
	"context"

	"github.com/Siposattila/gobkup/log"
	"github.com/disgoorg/disgo"
	discordgo "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

type discord struct {
	Client       webhook.Client
	WebhookId    string
	WebhookToken string
}

func NewDiscord(webhookId string, webhookToken string) AlertInterface {
	return &discord{WebhookId: webhookId, WebhookToken: webhookToken}
}

func (d *discord) Start() {
	log.GetLogger().Normal("Starting discord client...")
	log.GetLogger().Normal("Discord client version: (disgo)" + disgo.Version)

	d.Client = webhook.New(snowflake.MustParse(d.WebhookId), d.WebhookToken)
	log.GetLogger().Success("Discord client started! Ready to send alerts!")
}

func (d *discord) Stop() {
	if d.Client != nil {
		d.Client.Close(context.Background())
		log.GetLogger().Normal("Stopping discord client...")
	}
}

func (d *discord) Send(message string) {
	if _, alertError := d.Client.CreateMessage(
		discordgo.NewWebhookMessageCreateBuilder().
			SetContent(message).
			Build(),
	); alertError != nil {
		log.GetLogger().Error("There was an error during the creation and sending of a discord alert: " + alertError.Error())
	}
}
