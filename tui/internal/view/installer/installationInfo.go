package installer

func GetInstallationContent() string {
	return `
# Algorand Node UI Installer :smiley_cat:

You are about to install an Algorand Node. Congratulations!
If this is a mistake, and you already have a node installed,
you can exit this installer and run the UI with the -d flag.

When you are ready press **i** to start the installation wizard.
This will guide you through the installation process.

# What to expect
After running the wizard a node will be installed and started
and the node UI will transition to the main screen.

Two directories will be created:
* nodeui_algod_data
* nodeui_algod_bin

The data directory will contain the node data and configuration
while the bin directory will contain goal, algod, and other
binaries associated with the node.

# Node management
If you need to manage the node directly for any reason, you can
use the goal command. For example, to stop the node you can run:
` + "```sh" + `
nodeui_algod_bin/goal node stop -d nodeui_algod_data
` + "```" + `

To start it again you can run:
` + "```sh" + `
nodeui_algod_bin/goal node start -d nodeui_algod_data
` + "```" + `
`
}
