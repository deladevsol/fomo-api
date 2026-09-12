package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deladevsol/fomo-api/client"
	"github.com/deladevsol/fomo-api/internal/proxy"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if e := run(ctx, os.Args[1:]); e != nil && !errors.Is(e, flag.ErrHelp) && ctx.Err() == nil {
		log.Print(e)
		os.Exit(1)
	}
}
func run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		fmt.Println("fomo-api serve [--listen 127.0.0.1:8787]\nfomo-api get [--url URL] /v2/thesis\nfomo-api user [--url URL] HANDLE\nfomo-api resolve [--url URL] HANDLE")
		return nil
	}
	command := args[0]
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	base := fs.String("url", client.DefaultURL, "public API URL")
	listen := fs.String("listen", "127.0.0.1:8787", "local listen address")
	if e := fs.Parse(args[1:]); e != nil {
		return e
	}
	key := os.Getenv("FOMO_PUBLIC_API_KEY")
	if command == "serve" {
		h, e := proxy.New(*base, key)
		if e != nil {
			return e
		}
		s := &http.Server{Addr: *listen, Handler: h, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 90 * time.Second, MaxHeaderBytes: 64 << 10}
		go func() {
			<-ctx.Done()
			shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = s.Shutdown(shutdown)
		}()
		log.Printf("listening on %s", *listen)
		e = s.ListenAndServe()
		if e == http.ErrServerClosed {
			return nil
		}
		return e
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("%s requires one argument", command)
	}
	c, e := client.New(client.Config{BaseURL: *base, APIKey: key})
	if e != nil {
		return e
	}
	var out any
	switch command {
	case "get":
		var v json.RawMessage
		e = c.Request(ctx, http.MethodGet, fs.Arg(0), &v)
		out = v
	case "user":
		out, e = c.UserByHandle(ctx, fs.Arg(0))
	case "resolve":
		out, e = c.Resolve(ctx, fs.Arg(0))
	default:
		return fmt.Errorf("unknown command %q", command)
	}
	if e != nil {
		return e
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}
