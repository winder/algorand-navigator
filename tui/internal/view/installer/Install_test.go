package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallTest(t *testing.T) {
	rootDir := t.TempDir()

	data := path.Join(rootDir, "nodeui_algod_data")
	bin := path.Join(rootDir, "nodeui_algod_bin")
	err := os.Mkdir(data, 0755)
	require.NoError(t, err)
	err = os.Mkdir(bin, 0755)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(bin, "update.sh"), []byte(updateScript), 0755)
	require.NoError(t, err)

	c := exec.Command(fmt.Sprintf("%s/update.sh", bin), "-i", "-c", "stable", "-p", bin, "-d", data, "-i", "-g", "mainnet")
	out, err := c.CombinedOutput()
	fmt.Println(out)
	require.NoError(t, err)
}
