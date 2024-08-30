package commands

import (
	"context"
	"fmt"

	"os"
	"os/signal"

	"github.com/idreaminteractive/goreload/internal/hotreload"
	"github.com/urfave/cli/v2"
)

// Run the hot reload server
func Serve(cCtx *cli.Context, hostUrl string) error {
	_, port, err := hotreload.ValidateUrl(hostUrl)
	if err != nil {
		fmt.Fprintf(cCtx.App.ErrWriter, "Invalid host given @ %s: %v\n", hostUrl, err)
		return err
	}

	fmt.Fprintf(cCtx.App.Writer, "Hot reload server starting up on port %d\n", port)

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(cCtx.Context)
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; fmt.Fprint(cCtx.App.Writer, "Received SIGINT\n"); cancel() }()

	// init our program
	// sets it all up

	s := hotreload.InitHotReloadServer(port, hotreload.WithLogger(&cCtx.App.Writer))

	// // Execute program.
	if err := s.Run(ctx); err != nil {
		fmt.Fprint(cCtx.App.Writer, "Closing our hot reload due to startup error\n")
		s.Close()
		fmt.Fprintln(cCtx.App.ErrWriter, err)
		return err
	}

	// wait + receive on sigint
	<-ctx.Done()
	fmt.Fprint(cCtx.App.Writer, "Closing hot reload server\n")

	// clean up all the things
	if err := s.Close(); err != nil {
		fmt.Fprintln(cCtx.App.ErrWriter, err)
		return err
	}
	return nil
}
