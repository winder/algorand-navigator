// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package messages

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Requestor provides an opaque pointer for an algod client.
type Requestor struct {
	Client     *algod.Client
	url        string
	adminToken string
	token      string
	dataDir    string
}

// MakeRequestor builds the requestor object.
func MakeRequestor(url, token, adminToken, dataDir string) *Requestor {
	client, err := algod.MakeClient(url, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem creating client connection: %s\n", err.Error())
		os.Exit(1)
	}

	return &Requestor{
		url:        url,
		token:      token,
		adminToken: adminToken,
		Client:     client,
		dataDir:    dataDir,
	}
}

// NetworkMsg holds network information.
type NetworkMsg struct {
	GenesisID   string
	GenesisHash types.Digest
	NodeVersion string
	Err         error
}

func formatVersion(ver models.Version) string {
	return fmt.Sprintf("%s %d.%d.%d (%s)",
		ver.Build.Channel,
		ver.Build.Major,
		ver.Build.Major,
		ver.Build.BuildNumber,
		ver.Build.CommitHash)
}

// GetConfigs returns the node config.json file if possible.
func (r Requestor) GetConfigs() string {
	if r.dataDir == "" {
		return "data directory not set"
	}
	configs, err := os.ReadFile(path.Join(r.dataDir, "config.json"))
	if err != nil {
		return "config.json file not found"
	}
	return string(configs)
}

// GetNetworkCmd provides a tea.Cmd for fetching a NetworkMsg.
func (r Requestor) GetNetworkCmd() tea.Cmd {
	return func() tea.Msg {
		ver, err := r.Client.Versions().Do(context.Background())
		if err != nil {
			return NetworkMsg{
				Err: err,
			}
		}

		var digest types.Digest
		if len(ver.GenesisHash) != len(digest) {
			return NetworkMsg{
				Err: fmt.Errorf("unexpected genesis hash, wrong number of bytes"),
			}
		}
		copy(digest[:], ver.GenesisHash)

		return NetworkMsg{
			GenesisID:   ver.GenesisID,
			GenesisHash: digest,
			NodeVersion: formatVersion(ver),
		}
	}
}

// StatusMsg has node status information.
type StatusMsg struct {
	Status models.NodeStatus
	Error  error
}

// GetStatusCmd provides a tea.Cmd for fetching a StatusMsg.
func (r Requestor) GetStatusCmd() tea.Cmd {
	return func() tea.Msg {
		resp, err := r.Client.Status().Do(context.Background())
		//s, err := s.node.Status()
		return StatusMsg{
			Status: resp,
			Error:  err,
		}
	}
}

// AccountStatusMsg has account balance information.
type AccountStatusMsg struct {
	Balances map[types.Address]map[uint64]uint64
	Err      error
}

// GetAccountStatusCmd provides a tea.Cmd for fetching a AccountStatusMsg.
func (r Requestor) GetAccountStatusCmd(accounts []types.Address) tea.Cmd {
	return func() tea.Msg {
		var rval AccountStatusMsg
		rval.Balances = make(map[types.Address]map[uint64]uint64)

		for _, acct := range accounts {
			resp, err := r.Client.AccountInformation(acct.String()).Do(context.Background())
			if err != nil {
				return AccountStatusMsg{
					Err: err,
				}
			}
			rval.Balances[acct] = make(map[uint64]uint64)

			// algos at the special index
			rval.Balances[acct][0] = resp.Amount

			// everything else
			for _, holding := range resp.Assets {
				rval.Balances[acct][holding.AssetId] = holding.Amount
			}
		}

		return rval
	}
}

func doFastCatchupRequest(rootURL, adminToken, verb, network string) error {
	if adminToken == "" {
		return fmt.Errorf("cannot use fast catchup without an admin token")
	}
	resp, err := http.Get(fmt.Sprintf("https://algorand-catchpoints.s3.us-east-2.amazonaws.com/channel/%s/latest.catchpoint", network))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	catchpoint := strings.Replace(string(body), "#", "%23", 1)

	//start fast catchup
	url := fmt.Sprintf("%s/v2/catchup/%s", rootURL, catchpoint)
	url = url[:len(url)-1] // remove \n
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Algo-Api-Token", string(adminToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// StartFastCatchup attempts to start fast catchup for a given network.
func (r Requestor) StartFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
		err := doFastCatchupRequest(r.url, r.adminToken, http.MethodPost, network)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// StopFastCatchup attempts to stop fast catchup for a given network.
func (r Requestor) StopFastCatchup(network string) tea.Cmd {
	return func() tea.Msg {
		err := doFastCatchupRequest(r.url, r.adminToken, http.MethodDelete, network)
		if err != nil {
			panic(err)
		}
		return nil
	}
}
