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
	log "github.com/inconshreveable/log15"

	"github.com/netsec-ethz/scion/go/hps/conf"
	"github.com/netsec-ethz/scion/go/lib/common"
	"github.com/netsec-ethz/scion/go/lib/ctrl"
	"github.com/netsec-ethz/scion/go/lib/ctrl/path_mgmt"
	"github.com/netsec-ethz/scion/go/lib/snet"
)

var ()

type HPS struct {
	// Id is the SCION element ID, e.g. "hps4-21-9".
	Id string
	// confDir is the directory containing the configuration file.
	confDir string
	// dispPath is the path to the dispatcher
	dispPath string
	// sciondPath is the path to SCIOND
	sciondPath string
}

func NewHPS(id, confDir, dispPath, sciondPath string) (*HPS, *common.Error) {
	// TODO(shitz): add metrics
	hps := &HPS{
		Id:         id,
		confDir:    confDir,
		dispPath:   dispPath,
		sciondPath: sciondPath,
	}
	if cerr := hps.setup(); cerr != nil {
		return nil, cerr
	}
	return hps, nil
}

func (h *HPS) setup() *common.Error {
	var cerr *common.Error
	if cerr = conf.Load(h.confDir); cerr != nil {
		return cerr
	}
	return nil
}

func (h *HPS) Run() *common.Error {
	c := conf.Get()
	if cerr := snet.Init(c.IA, h.sciondPath, h.dispPath); cerr != nil {
		return common.NewError("Unable to initialize SCION network", "err", cerr)
	}
	log.Debug("SCION network successfully initialized")
	// TODO(shitz): Get local address from topo.
	conn, err := snet.ListenSCION("udp4", &local)
	if err != nil {
		return common.NewError("Unable to listen", "err", err)
	}
	log.Debug("Listening", "local", local)
	receive(conn)
	return nil
}

// receive is the main SCION packet handling loop.
func receive(conn *snet.Conn) {
	// Main loop
	readBuf := make(common.RawBytes, 1<<16)
	for {
		// Reset readBuf
		readBuf = readBuf[:cap(readBuf)]
		read, a, err := conn.ReadFromSCION(readBuf)
		if err != nil {
			log.Error("Error reading from conn", "conn", conn, "err", err)
			continue
		}
		// Slice readBuf
		readBuf = readBuf[:read]
		cpld, cerr := ctrl.NewPldFromRaw(readBuf)
		if cerr != nil {
			log.Error("Error parsing ctrl payload", "err", cerr)
			continue
		}
		// Determine the type of SCION control payload.
		u0, cerr := cpld.Union0()
		if cerr != nil {
			log.Error(cerr.String())
			continue
		}
		switch u0 := u0.(type) {
		case *path_mgmt.SegReq:
			log.Debug("Path request received!", "req", u0)
		case *path_mgmt.SegReg:
			log.Debug("Path registration received", "reg", u0)
		default:
			log.Error("Unsupported payload", "type", common.TypeOf(u0))
			continue
		}
	}
}
