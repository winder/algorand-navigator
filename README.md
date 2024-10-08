# Algorand Navigator

Terminal UI for remote Algorand node management.

![Example Screenshot](images/demo.png)

# Install
## Download
See the GitHub releases and download the binary for your platform.

## Source
Use go1.20.5 or later and build with `make`.

# Usage
With no options, the UI will be displayed instead of starting a service.

## Local Algod
```
~$ ALGORAND_DATA=path/to/data/dir ./algorand-navigator
```
## Remote Algod
```
~$ ./algorand-navigator -t <algod api token> -u http://<url>
```

# Run as a service

The preferred method for running the navigator UI is as a service running alongside algod. By passing a port using `-p` or `--tui-port` an SSH server is started and can host the UI for multiple clients.

A tool like [wishlist](https://github.com/charmbracelet/wishlist#wishlist) can be used to interactively select between multiple node deployments. In the screenshot below you can see a sample ssh config file, and the UI wishlist provides to select which navigator to connect to.

![Wishlist Example](images/wishlist_example.png)

# Features

## Status

Realtime node status, including detailed fast-catchup progress.

## Block Explorer

Display realtime block data, drill down into a block to see all of the transactions and transaction details.

## Utilities

Start a fast catchup with the press of a key, and more (if you build it)!

## Built in documentation

[Kind of](tui/internal/bubbles/about/help.go).

# Architecture

Built using [Bubble Tea](https://github.com/charmbracelet/bubbletea). Node information is collected from the Algod REST API using the [go SDK](https://github.com/algorand/go-algorand-sdk), and from reading files on disk.

Each box on the screen is a "bubble", they manage state independently with an event update loop. Events are passed to each bubble, which have the option of consuming the event and/or passing it along to any nested bubbles. When processing the event, they may optionally add follow-up tasks which the scheduling engine would execute asynchronously. Follow-up tasks may optionally create more events which would be processed in turn using the same mechanism.

When displaying the UI, each bubble is asked to renders itself and they are finally joined together for final rendering using [lipgloss](https://github.com/charmbracelet/lipgloss). Web development aficionado may recognize this pattern as [The Elm Architecture](https://guide.elm-lang.org/architecture/).

There are some quirks to this approach. The main one is that bubbletea is a rendering engine, NOT a window manager. This means that things like window heights and widths must be self-managed. Any mismanagement leads to very strange artifacts as the rendering engine tries to fit too many, or too few lines to a fixed sized terminal.

# Contributing

Contributions are welcome! There are no plans to actively maintain this project, so if you find it useful please consider helping out.

# How to create a new release

1. Create a tag: `git tag -a v_._._ -m "v_._._" && git push origin v_._._`
2. Push the tag.
3. CI should create a release, attach it to GitHub and publish images to docker hub.
