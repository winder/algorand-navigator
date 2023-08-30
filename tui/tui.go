package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"

	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/tui/internal/constants"
	"github.com/algorand/node-ui/tui/internal/view/setup"
)

func getTeaHandler(model tea.Model) bm.Handler {
	return func(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
		return model, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
	}
}

// Start ...
func Start(args args.Arguments) {
	model := setup.New(args)

	// Run directly
	if args.TuiPort == 0 {
		p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error in UI: %v", err)
			os.Exit(1)
		}

		fmt.Printf("\nUI Terminated, shutting down node.\n")
		os.Exit(0)
	}

	// Run on ssh server.
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	sshServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", constants.Host, args.TuiPort)),
		wish.WithHostKeyPath(path.Join(dirname, ".ssh/term_info_ed25519")),
		wish.WithMiddleware(
			bm.Middleware(getTeaHandler(model)),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", constants.Host, args.TuiPort)
	go func() {
		if err = sshServer.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := sshServer.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
