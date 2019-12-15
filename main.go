package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blakey22/line-notify-gateway/pkg/flag"
	"github.com/blakey22/line-notify-gateway/pkg/handler"
	_ "github.com/blakey22/line-notify-gateway/pkg/handler/all"
	"github.com/blakey22/line-notify-gateway/pkg/line"
	"github.com/jessevdk/go-flags"
	"golang.org/x/sync/errgroup"
)

func main() {
	_, err := flags.ParseArgs(&flag.Options, os.Args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(1)
		} else {
			panic(err)
		}
	}

	if handler.Count() == 0 {
		panic("no handler has been registered")
	}

	if len(flag.Options.Secret) == 0 {
		log.Print("WARNING! no secret token is set, your gateway might expose to everyone")
	} else if len(flag.Options.Token) == 0 {
		log.Print("Token of LINE notify is not set")
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", flag.Options.Host, flag.Options.Port)
	log.Printf("Line Notification Gateway starts at %s", addr)

	g, ctx := errgroup.WithContext(context.Background())
	sendc := make(chan line.Notification, handler.Count()*2)
	defer close(sendc)

	notifier := line.NewNotifier(flag.Options.Endpoint, flag.Options.Token)
	g.Go(func() error {
		return notifier.Run(ctx, sendc)
	})

	handler.Setup(sendc)
	g.Go(func() error {
		return http.ListenAndServe(addr, nil)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("gateway is stopped due to an error: %+v", err)
	}
}
