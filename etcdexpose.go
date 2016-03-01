package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/coreos/etcd/client"

	"github.com/upfluence/etcdexpose/handler"
	"github.com/upfluence/etcdexpose/handler/multiple"
	"github.com/upfluence/etcdexpose/handler/single"
	"github.com/upfluence/etcdexpose/runner"
	"github.com/upfluence/etcdexpose/utils"
	"github.com/upfluence/etcdexpose/watcher"
	"github.com/upfluence/etcdexpose/watcher/etcd"
	time_watcher "github.com/upfluence/etcdexpose/watcher/time"
)

const (
	currentVersion = "0.1.0-Meh"
	bufferSize     = 50
)

var (
	flagset = flag.NewFlagSet("etcdexpose", flag.ExitOnError)
	flags   = struct {
		Version    bool
		Multiple   bool
		Server     string
		Template   string
		Namespace  string
		HealthPath string
		Key        string
		Interval   uint
		Retry      uint
		RetryDelay time.Duration
		Timeout    time.Duration
		Ttl        time.Duration
		Port       uint
		CheckPort  uint
	}{}
)

func usage() {
	fmt.Fprintf(os.Stderr, `
  NAME
  etcdexpose - expose clusters from discovery directories

  USAGE
  etcdexpose [options]

  OPTIONS
  `)
	flagset.PrintDefaults()
}

func init() {
	flagset.BoolVar(&flags.Version, "version", false, "Print the version and exit")
	flagset.BoolVar(&flags.Version, "v", false, "Print the version and exit")

	flagset.BoolVar(&flags.Multiple, "multiple", false, "Expose all registered keys")
	flagset.BoolVar(&flags.Multiple, "m", false, "Expose all registered keys")

	flagset.StringVar(&flags.Server, "server", "http://127.0.0.1:4001", "Location of the etcd server")
	flagset.StringVar(&flags.Server, "s", "http://127.0.0.1:4001", "Location of the etcd server")

	flagset.StringVar(&flags.Template, "template", "http://{{.Value}}:{{.Port}}", "Template to apply")
	flagset.StringVar(&flags.Template, "t", "http://{{.Value}}:{{.Port}}", "Template to apply")

	flagset.StringVar(&flags.Namespace, "namespace", "/", "Discovery directory to watch")
	flagset.StringVar(&flags.Namespace, "n", "/", "Discovery directory to watch")

	flagset.StringVar(&flags.Key, "key", "/key", "Key to expose")
	flagset.StringVar(&flags.Key, "k", "/key", "key to expose")

	flagset.StringVar(&flags.HealthPath, "health-check", "/", "Path to use to perform healthCheck")
	flagset.StringVar(&flags.HealthPath, "h", "/", "Path to use to perform healthCheck")

	flagset.UintVar(&flags.Interval, "interval", 0, "Perform an update at regular interval if > 0")
	flagset.UintVar(&flags.Interval, "i", 0, "Perform an update at regulat interfal if > 0")

	flagset.UintVar(&flags.Retry, "retry", 0, "Healthcheck retry")
	flagset.UintVar(&flags.Retry, "r", 0, "Healtcheck retry")
	flagset.DurationVar(&flags.RetryDelay, "retry-delay", 5*time.Second, "Healtcheck retry delay")
	flagset.DurationVar(&flags.RetryDelay, "rd", 5*time.Second, "Healtcheck retry delay")

	flagset.DurationVar(&flags.Timeout, "client-timeout", 5*time.Second, "Client timeout")
	flagset.DurationVar(&flags.Timeout, "ct", 5*time.Second, "Client timeout")

	flagset.DurationVar(&flags.Ttl, "ttl", 0, "Key time to live")

	flagset.UintVar(&flags.Port, "port", 0, "Port to expose")
	flagset.UintVar(&flags.Port, "p", 0, "Port to expose")
	flagset.UintVar(&flags.CheckPort, "check-port", 0, "Check port to use")
}

func main() {
	flagset.Parse(os.Args[1:])
	flagset.Usage = usage

	if len(os.Args) < 2 {
		flagset.Usage()
		os.Exit(0)
	}

	if flags.Version {
		fmt.Printf("etcdexpose v%s", currentVersion)
		os.Exit(0)
	}

	if flags.Port == 0 {
		fmt.Println("You must provide a valid port to expose with -p flag")
		os.Exit(1)
	}

	if flags.CheckPort == 0 {
		flags.CheckPort = flags.Port
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	cfg := client.Config{
		Endpoints:               []string{flags.Server},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	renderer, err := utils.NewValueRenderer(flags.Template, flags.Port)

	if err != nil {
		log.Fatalf("Invalid template given")
	}

	healthCheck, err := utils.NewHealthCheck(
		flags.HealthPath,
		flags.CheckPort,
		flags.Retry,
		flags.RetryDelay,
		flags.Timeout,
	)

	if err != nil {
		log.Fatalf("Invalid format given")
	}

	namespace_client := utils.NewEtcdClient(
		kapi,
		flags.Namespace,
		flags.Key,
		flags.Ttl,
	)

	var handler handler.Handler = nil

	if flags.Multiple {
		handler = multiple.NewMutlipleValueExpose(
			namespace_client,
			renderer,
			healthCheck,
		)

	} else {
		handler = single.NewSingleValueExpose(
			namespace_client,
			renderer,
			healthCheck,
		)
	}

	watchers := []watcher.Watcher{
		etcd.NewWatcher(kapi, flags.Namespace, bufferSize),
	}

	if flags.Interval > 0 {
		timeWatcher := time_watcher.NewWatcher(
			time.Duration(flags.Interval)*time.Second,
			bufferSize,
		)
		watchers = append(watchers, timeWatcher)

	}

	runner := runner.NewRunner(
		handler,
		watchers,
		len(watchers)*bufferSize,
	)

	go func() {
		s := <-sigch
		log.Printf("Received signal [%v] stopping application\n", s)
		runner.Stop()
		os.Exit(0)
	}()

	for {
		log.Println("Starting runner...")
		runner.Start()
		log.Println("Runner exited, Stopping...")
		runner.Stop()
		log.Println("waiting 5s before retry ...")
		time.Sleep(5 * time.Second)
	}
}
