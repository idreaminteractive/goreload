package commands

import (
	"fmt"
	"net/http"

	"github.com/idreaminteractive/goreload/internal/hotreload"
)

// Send the hot reload signal
func SignalReload(url string) error {
	host, _, err := hotreload.ValidateUrl(url)
	if err != nil {
		return err
	}
	// ping to our post
	pingUrl := host + "/hotreload"

	hc := http.Client{}
	req, err := http.NewRequest("POST", pingUrl, nil)
	if err != nil {
		panic(err)
	}
	_, err = hc.Do(req)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
	return nil
}
