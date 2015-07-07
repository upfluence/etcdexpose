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
	"os"
	"os/signal"

	"github.com/upfluence/etcdexpose/etcdexpose"
)

const currentVersion = "0.0.1"

var (
	flagset = flag.NewFlagSet("etcdexpose", flag.ExitOnError)
	flags   = struct {
		Version   bool
		All       bool
		Server    string
		Template  string
		Namespace string
		Key       string
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

	flagset.BoolVar(&flags.All, "all", false, "Expose all registered keys")
	flagset.BoolVar(&flags.All, "a", false, "Expose all registered keys")

	flagset.StringVar(&flags.Server, "server", "http://127.0.0.1:4001", "Location of the etcd server")
	flagset.StringVar(&flags.Server, "s", "http://127.0.0.1:4001", "Location of the etcd server")

	flagset.StringVar(&flags.Template, "template", "http://{}", "Template to apply")
	flagset.StringVar(&flags.Template, "t", "http://{}", "Template to apply")

	flagset.StringVar(&flags.Namespace, "namespace", "/", "Discovery directory to watch")
	flagset.StringVar(&flags.Namespace, "n", "/", "Discovery directory to watch")

	flagset.StringVar(&flags.Key, "key", "/key", "Key to expose")
	flagset.StringVar(&flags.Key, "k", "/key", "key to expose")
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

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		<-sigch
		os.Exit(0)
	}()

	watcher := etcdexpose.NewEtcdWatcher(
		flags.Namespace,
		[]string{flags.Server})

	go watcher.Start()

	for {
		select {
		case event := <-watcher.EventChan:
			fmt.Printf("%s", event.Action)
		case err := <-watcher.ErrorChan:
			fmt.Printf("Error %s", err)
		}
	}

}
