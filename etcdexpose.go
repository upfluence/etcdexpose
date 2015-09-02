/*
   Copyright 2014 Upfluence, Inc.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/upfluence/etcdexpose/etcdexpose"
)

const currentVersion = "0.0.7"

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
		Ttl        uint64
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

	flagset.Uint64Var(&flags.Ttl, "ttl", 0, "Key time to live")

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

	client := etcd.NewClient([]string{flags.Server})

	renderer, err := etcdexpose.NewValueRenderer(flags.Template, flags.Port)

	if err != nil {
		log.Fatalf("Invalid template given")
	}

	healthCheck := etcdexpose.NewHealthCheck(
		flags.HealthPath, flags.CheckPort,
	)

	namespace_client := etcdexpose.NewEtcdClient(
		client,
		flags.Namespace,
		flags.Key,
		flags.Ttl,
	)

	etcdWatcher := etcdexpose.NewEtcdWatcher(
		flags.Namespace,
		client,
	)

	var handler etcdexpose.Handler = nil

	if flags.Multiple {
		handler = etcdexpose.NewMutlipleValueExpose(
			namespace_client,
			renderer,
			healthCheck,
		)

	} else {
		handler = etcdexpose.NewSingleValueExpose(
			namespace_client,
			renderer,
			healthCheck,
		)
	}

	runner := etcdexpose.NewRunner(handler)
	runner.AddWatcher(etcdWatcher)

	if flags.Interval > 0 {
		timeWatcher := etcdexpose.NewTimeWatcher(
			time.Duration(flags.Interval),
			time.Second,
		)

		runner.AddWatcher(timeWatcher)
	}

	go func() {
		<-sigch
		runner.Stop()
		os.Exit(0)
	}()

	for {
		runner.Start()
		runner.Stop()
		log.Printf("Runner exited, waiting 5s ...")
		time.Sleep(5 * time.Second)
	}
}
