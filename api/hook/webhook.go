package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/httpcontrol"
	"github.com/ovh/tat"
	"github.com/spf13/viper"
)

func sendWebHook(hook tat.Hook, path string) error {

	data, err := json.Marshal(hook)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", path, bytes.NewReader(data))

	c := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: 5 * time.Second,
			MaxTries:       3,
		},
	}

	resp, err := c.Do(req)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return fmt.Errorf("sendWebHook >> err:%s", err)
	}

	if !viper.GetBool("production") {
		body, errb := ioutil.ReadAll(resp.Body)
		if errb != nil {
			return fmt.Errorf("sendWebHook >> Error with ioutil.ReadAll %s", errb.Error())
		}
		log.Debugf("Response from webhook %s", body)
	}

	return nil
}
