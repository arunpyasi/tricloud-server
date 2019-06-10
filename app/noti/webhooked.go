package noti

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/logg"
)

type WebHooked struct {
}

func (w *WebHooked) Init()           {}
func (w *WebHooked) GetName() string { return "webhooked" }
func (w *WebHooked) PushNotification(ctx context.Context, token string, data map[string]string) error {

	var outbody = map[string]string{
		"content": mapToStr(data),
	}

	outbyte, err := json.Marshal(outbody)
	if err != nil {
		logg.Info("Encoding Error")
		return err
	}

	req, err := http.NewRequest("POST", token, bytes.NewBuffer(outbyte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logg.Info("WEBHOOKED ERROR")
		logg.Info(err)
		return err
	}
	defer resp.Body.Close()

	logg.Info("response Status:")
	logg.Info(resp.Status)
	logg.Info("response Headers:")
	logg.Info(resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logg.Info("response Body:")
	logg.Info(string(body))

	return nil
}

func mapToStr(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
