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

// Package conf holds all of the global HPS state, for access by the
// HPS' various packages.
package conf

import (
	"path/filepath"
	"sync"

	"github.com/netsec-ethz/scion/go/lib/addr"
	"github.com/netsec-ethz/scion/go/lib/common"
	"github.com/netsec-ethz/scion/go/lib/topology"
)

type Conf struct {
	// Topo contains the names of all local infrastructure elements, a map
	// of interface IDs to routers, and the actual topology.
	Topo *topology.Topo
	// IA is the current ISD-AS.
	IA *addr.ISD_AS
	// Dir is the configuration directory.
	Dir string
}

// Load sets up the configuration, loading it from the supplied config directory.
func Load(confDir string) *common.Error {
	var cerr *common.Error
	// Declare a new Conf instance, and load the topology config.
	c := &Conf{}
	c.Dir = confDir
	topoPath := filepath.Join(c.Dir, topology.CfgName)
	if c.Topo, cerr = topology.LoadFromFile(topoPath); cerr != nil {
		return cerr
	}
	c.IA = c.Topo.ISD_AS
	Set(c)
	return nil
}

var c *Conf

var lock sync.RWMutex

func Get() *Conf {
	lock.RLock()
	defer lock.RUnlock()
	return c
}

func Set(newConf *Conf) {
	lock.Lock()
	defer lock.Unlock()
	c = newConf
}
