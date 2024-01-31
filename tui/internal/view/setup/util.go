package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/types"

	"github.com/algorand/node-ui/messages"
)

func getRequestor(algodDataDir, url, token, adminToken string) (*messages.Requestor, error) {
	// Initialize from -d, ALGORAND_DATA, or provided URL/Token

	if algodDataDir != "" && (url != "" || token != "") {
		algodDataDir = ""
		fmt.Println("ignoring ALGORAND_DATA/-d in favor of -u/-t")
	}

	// If url/token are missing, attempt to use environment variable.
	if algodDataDir != "" {
		netpath := filepath.Join(algodDataDir, "algod.net")
		tokenpath := filepath.Join(algodDataDir, "algod.token")
		adminTokenpath := filepath.Join(algodDataDir, "algod.admin.token")

		var netaddrbytes []byte
		netaddrbytes, err := os.ReadFile(netpath)
		if err != nil {
			return nil, fmt.Errorf("Unable to read URL from file (%s): %s\n", netpath, err.Error())
		}
		url = strings.TrimSpace(string(netaddrbytes))

		tokenBytes, err := os.ReadFile(tokenpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read Token from file (%s): %s\n", tokenpath, err.Error())
			os.Exit(1)
		}
		token = string(tokenBytes)

		adminTokenBytes, err := os.ReadFile(adminTokenpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read Token from file (%s): %s\n", adminTokenpath, err.Error())
			os.Exit(1)
		}
		adminToken = string(adminTokenBytes)
	}

	if url == "" || token == "" {
		return nil, fmt.Errorf("must provide a way to get the algod REST API")
	}

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	return messages.MakeRequestor(url, token, adminToken, algodDataDir), nil
}

func getRequestorOrExit(algodDataDir, url, token, adminToken string) *messages.Requestor {
	requestor, err := getRequestor(algodDataDir, url, token, adminToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem creating requestor: %s.\n", err.Error())
		os.Exit(1)
	}
	return requestor
}

func getAddressesOrExit(addrs []string) (result []types.Address) {
	failed := false
	for _, addr := range addrs {
		converted, err := types.DecodeAddress(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode address '%s': %s\n", addr, err.Error())
			failed = true
		}
		result = append(result, converted)
	}

	if failed {
		os.Exit(1)
	}

	return result
}
