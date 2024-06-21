package main

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/idreaminteractive/go-reload/internal/hotreload"
)

func ReloadComponent(hostUrl string) templ.Component {
	host, _, err := hotreload.ValidateUrl(hostUrl)
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
