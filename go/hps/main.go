// Copyright 2017 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"os"

	log "github.com/inconshreveable/log15"
	liblog "github.com/netsec-ethz/scion/go/lib/log"
)

var (
	local      snet.Addr
	id      = flag.String("id", "", "Element ID (Required. E.g. 'hps4-21-9')")
	confDir = flag.String("confd", ".", "Configuration directory")
	sciond     = flag.String("sciond", "", "Path to sciond socket")
	dispatcher = flag.String("dispatcher", "/run/shm/dispatcher/default.sock",
		"Path to dispatcher socket")
)

func init() {
	flag.Var((*Address)(&local), "local", "(Mandatory) address to listen on")
}

func main() {
	flag.Parse()
	if *id == "" {
		log.Crit("No element ID specified")
		os.Exit(1)
	}
	if *sciond == "" {
		*sciond = GetDefaultSCIONDPath(local.IA)
	}

	liblog.Setup(*id)
	defer liblog.LogPanicAndExit()

	hps, cerr := NewHPS(*id, *confDir, *dispatcher, *sciond)
	if cerr != nil {
		log.Crit("Startup failed", "err", cerr)
		liblog.Flush()
		os.Exit(1)
	}
	log.Debug("HPS started", "id", *id, "pid", os.Getpid())
	if cerr := hps.Run(); err != nil {
		log.Crit("Run failed", "err", cerr)
		liblog.Flush()
		os.Exit(1)
	}
}
