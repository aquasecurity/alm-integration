package outputs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aquasecurity/postee/data"
	"github.com/aquasecurity/postee/formatting"
	"github.com/aquasecurity/postee/layout"
)

const defaultSizeLimit = 10000

type SplunkOutput struct {
	Name         string
	Url          string
	Token        string
	EventLimit   int
	splunkLayout layout.LayoutProvider
}

func (splunk *SplunkOutput) GetName() string {
	return splunk.Name
}

func (splunk *SplunkOutput) Init() error {
	splunk.splunkLayout = new(formatting.HtmlProvider)
	log.Printf("Starting Splunk output %q....", splunk.Name)
	return nil
}

func (splunk *SplunkOutput) Send(d map[string]string) error {
	log.Printf("Sending a message to %q", splunk.Name)

	if splunk.EventLimit == 0 {
		splunk.EventLimit = defaultSizeLimit
	}
	if splunk.EventLimit < defaultSizeLimit {
		log.Printf("[WARNING] %q has a short limit %d (default %d)",
			splunk.Name, splunk.EventLimit, defaultSizeLimit)
	}

	if !strings.HasSuffix(splunk.Url, "/") {
		splunk.Url += "/"
	}

	scanInfo := new(data.ScanImageInfo)
	err := json.Unmarshal([]byte(d["src"]), scanInfo)
	if err != nil {
		log.Printf("sending to %q error: %v", splunk.Name, err)
		return err
	}

	eventFormat := "{\"sourcetype\": \"_json\", \"event\": "
	constLimit := len(eventFormat) - 1

	var fields []byte

	for {
		fields, err = json.Marshal(scanInfo)
		if err != nil {
			log.Printf("sending to %q error: %v", splunk.Name, err)
			return err
		}
		if len(fields) < splunk.EventLimit-constLimit {
			break
		}
		switch {
		case len(scanInfo.Resources) > 0:
			scanInfo.Resources = nil
			continue
		case len(scanInfo.Malwares) > 0:
			scanInfo.Malwares = nil
			continue
		case len(scanInfo.SensitiveData) > 0:
			scanInfo.SensitiveData = nil
			continue
		default:
			msg := fmt.Sprintf("Scan result for %q is large for %q , its size if %d (limit %d)",
				scanInfo.Image, splunk.Name, len(fields), splunk.EventLimit)
			log.Print(msg)
			return errors.New(msg)
		}
	}

	var buff bytes.Buffer
	buff.WriteString(eventFormat)
	buff.Write(fields)
	buff.WriteByte('}')

	req, err := http.NewRequest("POST", splunk.Url+"services/collector", &buff)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Splunk "+splunk.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Splunk sending error: failed response status %q. Body: %q", resp.Status, string(b))
		return errors.New("failed response status for Splunk sending")
	}
	log.Printf("Sending a message to %q was successful!", splunk.Name)
	return nil
}

func (splunk *SplunkOutput) Terminate() error {
	log.Printf("Splunk output %q terminated", splunk.Name)
	return nil
}

func (splunk *SplunkOutput) GetLayoutProvider() layout.LayoutProvider {
	return splunk.splunkLayout
}
