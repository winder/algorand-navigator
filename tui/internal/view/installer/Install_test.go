package installer

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path"
	"testing"
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

	command := fmt.Sprintf(`%s/update.sh -i -c "%s" -p "%s" -d "%s" -i -g "%s"`, t.TempDir(), "stable", bin, data, "stable")
	c := exec.Command(command)
	out, err := c.CombinedOutput()
	fmt.Println(out)
	require.NoError(t, err)
}
