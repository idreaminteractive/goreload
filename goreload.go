package goreload

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/idreaminteractive/goreload/internal/commands"
	"github.com/idreaminteractive/goreload/internal/hotreload"
)

// Programmatically trigger the hot reload for the server running @ reloadServerUrl
func SendReloadSignal(reloadServerUrl string) error {
	return commands.SignalReload(reloadServerUrl)
}

// Embed the JS to connect to the hot reload server and wait for SSE events @ reloadServerUrl
func ReloadComponent(reloadServerUrl string) templ.Component {
	host, _, err := hotreload.ValidateUrl(reloadServerUrl)
	if err != nil {
		panic(err)
	}
	url := host + "/hotreload"
	output := fmt.Sprintf(`
	<script type="text/javascript">

	(function () {
		let reloadSrc = window.goreload_reloadSrc || new EventSource("%s");
		reloadSrc.onmessage = (event) => {
		  if (event && event.data === "reload") {
			window.location.reload();
		  }
		};
		window.reloadSrc = reloadSrc;
	  })();

	</script>
	`, url)

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, output)
		return err

	})
}
