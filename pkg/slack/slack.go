package slack

import (
	"log"

	"github.com/slack-go/slack"
)

type SlackMessageInfo struct {
	Api       slack.Client
	ChannelID string
	UserName  string
	IconEmoji string
	Message   string
}

func NewSlackMessageInfo(token, channelId, userName, iconEmoji, message string) *SlackMessageInfo {
	return &SlackMessageInfo{
		Api:       *slack.New(token),
		ChannelID: channelId,
		UserName:  userName,
		IconEmoji: iconEmoji,
		Message:   message,
	}
}

func (i *SlackMessageInfo) PostMessage() error {
	if _, _, err := i.Api.PostMessage(
		i.ChannelID,
		slack.MsgOptionText(i.Message, false),
		slack.MsgOptionIconEmoji(i.IconEmoji),
		slack.MsgOptionUsername(i.UserName),
		slack.MsgOptionAsUser(false),
	); err != nil {
		return err
	}
	log.Printf("[INFO] Post message %v", i)
	return nil
}
