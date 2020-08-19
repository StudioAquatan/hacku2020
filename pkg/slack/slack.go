package slack

import "github.com/slack-go/slack"

type MessageInfo struct {
	Api       slack.Client
	ChannelID string
	UserName  string
	IconEmoji string
}

func NewMessageInfo(token, channelId, userName, iconEmoji string) *MessageInfo {
	return &MessageInfo{
		Api:       *slack.New(token),
		ChannelID: channelId,
		UserName:  userName,
		IconEmoji: iconEmoji,
	}
}

func (i *MessageInfo) PostMessage(t string) error {
	if _, _, err := i.Api.PostMessage(
		i.ChannelID,
		slack.MsgOptionText(t, false),
		slack.MsgOptionIconEmoji(i.IconEmoji),
		slack.MsgOptionUsername(i.UserName),
		slack.MsgOptionAsUser(false),
	); err != nil {
		return err
	}

	return nil
}
