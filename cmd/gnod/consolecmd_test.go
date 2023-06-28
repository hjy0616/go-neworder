// Copyright 2016 The go-neworder Authors
// This file is part of go-neworder.
//
// go-neworder is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-neworder is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-neworder. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/rand"
	"math/big"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/neworder-network/go-neworder/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 engine:1.0 eth:1.0 ethash:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 net:1.0 rpc:1.0 web3:1.0"
)

// spawns gnod with the given command line args, using a set of flags to minimise
// memory and disk IO. If the args don't set --datadir, the
// child g gets a temporary data directory.
func runMinimalGnod(t *testing.T, args ...string) *testgnod {
	// --ropsten to make the 'writing genesis to disk' faster (no accounts)
	// --networkid=1337 to avoid cache bump
	// --syncmode=full to avoid allocating fast sync bloom
	allArgs := []string{"--ropsten", "--networkid", "1337", "--authrpc.port", "0", "--syncmode=full", "--port", "0",
		"--nat", "none", "--nodiscover", "--maxpeers", "0", "--cache", "64",
		"--datadir.minfreedisk", "0"}
	return runGnod(t, append(allArgs, args...)...)
}

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

	// Start a gnod console, make sure it's cleaned up and terminate the console
	gnod := runMinimalGnod(t, "--miner.etherbase", coinbase, "console")

	// Gather all the infos the welcome message needs to contain
	gnod.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	gnod.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	gnod.SetTemplateFunc("gover", runtime.Version)
	gnod.SetTemplateFunc("gnodver", func() string { return params.VersionWithCommit("", "") })
	gnod.SetTemplateFunc("niltime", func() string {
		return time.Unix(0, 0).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")
	})
	gnod.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	gnod.Expect(`
Welcome to the Gnod JavaScript console!

instance: Gnod/v{{gnodver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{.Etherbase}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

To exit, press ctrl-d or type exit
> {{.InputLine "exit"}}
`)
	gnod.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestAttachWelcome(t *testing.T) {
	var (
		ipc      string
		httpPort string
		wsPort   string
	)
	// Configure the instance for IPC attachment
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\gnod` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ipc = filepath.Join(t.TempDir(), "gnod.ipc")
	}
	// And HTTP + WS attachment
	p := trulyRandInt(1024, 65533) // Yeah, sometimes this will fail, sorry :P
	httpPort = strconv.Itoa(p)
	wsPort = strconv.Itoa(p + 1)
	gnod := runMinimalGnod(t, "--miner.etherbase", "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182",
		"--ipcpath", ipc,
		"--http", "--http.port", httpPort,
		"--ws", "--ws.port", wsPort)
	t.Run("ipc", func(t *testing.T) {
		waitForEndpoint(t, ipc, 3*time.Second)
		testAttachWelcome(t, gnod, "ipc:"+ipc, ipcAPIs)
	})
	t.Run("http", func(t *testing.T) {
		endpoint := "http://127.0.0.1:" + httpPort
		waitForEndpoint(t, endpoint, 3*time.Second)
		testAttachWelcome(t, gnod, endpoint, httpAPIs)
	})
	t.Run("ws", func(t *testing.T) {
		endpoint := "ws://127.0.0.1:" + wsPort
		waitForEndpoint(t, endpoint, 3*time.Second)
		testAttachWelcome(t, gnod, endpoint, httpAPIs)
	})
	gnod.ExpectExit()
}

func testAttachWelcome(t *testing.T, gnod *testgnod, endpoint, apis string) {
	// Attach to a running gnod note and terminate immediately
	attach := runGnod(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("gnodver", func() string { return params.VersionWithCommit("", "") })
	attach.SetTemplateFunc("etherbase", func() string { return gnod.Etherbase })
	attach.SetTemplateFunc("niltime", func() string {
		return time.Unix(0, 0).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")
	})
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return gnod.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the Gnod JavaScript console!

instance: Gnod/v{{gnodver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{etherbase}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

To exit, press ctrl-d or type exit
> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running concurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}