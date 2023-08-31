package installer

func GetInstallationContent() string {
	return `
# Algorand Node UI Installer :smiley_cat:

You are about to install an Algorand Node. Congratulations!
If this is a mistake, and you already have a node you'd like
to connect to please exit and run the UI again using the -d
or -u/-t options.

When you are ready press **i** to start the installation wizard.
This will guide you through the installation process.

# What to expect
The wizard will ask you to select a network to install. You can
install multiple networks, but you will need to run the wizard
for each network you want to install.

If a network already exists, the wizard will check for an update
and then start or restart the node.

Once the node is running the UI will automatically connect to it.

# Node management
The node will not be shutdown when you exit the UI. You can
manage the node using the goal command line tool. It will be
located in your user config directory. More information about
this is provided at the end of the installation.
`
}
