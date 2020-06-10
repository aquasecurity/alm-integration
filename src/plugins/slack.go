package plugins

import (
	"bytes"
	"data"
	"encoding/json"
	"formatting"
	"layout"
	"log"
	"settings"
	"strings"

	slackAPI "slack-api"
)

const (
	slackBlockLimit = 49
)

type SlackPlugin struct {
	Url           string
	SlackSettings *settings.Settings
	slackLayout   layout.LayoutProvider
}

func (slack *SlackPlugin) Init() error {
	slack.slackLayout = new(formatting.SlackMrkdwnProvider)
	log.Printf("Starting Slack plugin %q....", slack.SlackSettings.PluginName)
	return nil
}

func clearSlackText(text string) string {
	s := strings.ReplaceAll(text, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func buildSlackBlock(title string, data []byte) []byte {
	var content bytes.Buffer
	content.WriteByte('{')
	content.WriteString("\"blocks\":")
	content.WriteByte('[')
	content.WriteString(title)
	content.Write(data)
	content.WriteByte(']')
	content.WriteByte('}')
	return content.Bytes()
}

func (slack *SlackPlugin) Send(input map[string]string) error {
	log.Printf("Sending via Slack %q", slack.SlackSettings.PluginName)
	title := clearSlackText(slack.slackLayout.TitleH2(input["title"]))
	var body string
	if strings.HasSuffix(input["description"], ",") {
		body = strings.TrimSuffix(input["description"], ",")
	} else {
		body = input["description"]
	}
	body = "[" + clearSlackText(body) + "]"
	rawBlock := make([]data.SlackBlock, 0)
	err := json.Unmarshal([]byte(body), &rawBlock)
	if err != nil {
		log.Printf("Unmarshal slack sending error: %v", err)
		return err
	}

	length := len(rawBlock)

	if length >= slackBlockLimit {
		message := buildShortMessage(slack.SlackSettings.AquaServer, input["url"], slack.slackLayout)
		if err := slackAPI.SendToUrl(slack.Url, buildSlackBlock(title, []byte(message))); err != nil {
			log.Printf("Slack Sending Error: %v", err)
		}
		log.Printf("Sending via Slack %q was successful!", slack.SlackSettings.PluginName)
	} else {
		for n := 0; n < length; {
			d := length - n
			if d >= 49 {
				d = 49
			}
			cutData, _ := json.Marshal(rawBlock[n : n+d])
			cutData = cutData[1 : len(cutData)-1]
			if err := slackAPI.SendToUrl(slack.Url, buildSlackBlock(title, cutData)); err != nil {
				log.Printf("Slack Sending Error: %v", err)
			} else {
				log.Printf("Sending [%d/%d part] to %q was successful!",
					int(n/49)+1, int(length/49)+1,
					slack.SlackSettings.PluginName)
			}
			n += d
		}
	}
	return nil
}

func (slack *SlackPlugin) Terminate() error {
	log.Printf("Slack plugin %q terminated", slack.SlackSettings.PluginName)
	return nil
}

func (slack *SlackPlugin) GetLayoutProvider() layout.LayoutProvider {
	return slack.slackLayout
}

func (slack *SlackPlugin) GetSettings() *settings.Settings {
	return slack.SlackSettings
}
