package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "io"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/mmcdole/gofeed"
	"github.com/muesli/termenv"

	// "github.com/charmbracelet/lipgloss"

	f "github.com/Raven-bit/gorssh/cups"
)

// TODO: Stop using string append and start using whatever the String
// Buffer thing is; because golang strings are immutable and
// concatenation is SLOW.

// TODO: It looks like every field is optional? Actually go find out
// what Atom/RSS specs say here. But, assume every field may not be
// present.

// TODO: If content is not present, go fetch the content from the URL
// when we are asked to view the page. Oh, and all the network I/O
// should be done through tea Cmds, not just us.

// TODO: Add viper so we can have options to add items.

// TODO: Add a top level menu so we can decide if we're manipulating the
// list of feeds we know about, viewing entries on a feed, viewing
// entries on ALL feeds.

// TODO: Implement some kind of caching/file store so we keep feed
// entries we know about and can append to them. (Do RSS feeds support
// some kind of short changed? call? Can we use it? Do they support
// pagination? Time ranges? Or do we just.. get whatever we get, and
// we'd better have saved all the old stuff ourselves?)

// TODO: Make it an SSH service for extra cool points.

// TODO: We should make all of our output Markdown and just leverage
// glamour hard - maybe some custom themes just because we can?

// TODO: We've got an FR we can partially satisfy by offering a
// focus-on-feed option

const (
	host = "localhost"
	port = 23234
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			teaMiddleWare(),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

func teaHandler(s ssh.Session) *tea.Program {
	_, _, active := s.Pty()
	if !active {
		fmt.Println("no active terminal, skipping")
		_ = s.Exit(1)
		return nil
	}
	p := tea.NewProgram(
		f.NewFeedModel(
			getEntries("https://www.phoronix.com/rss.php"),
		),
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(),
	)

	return p
}

func teaMiddleWare() wish.Middleware {
	return bm.MiddlewareWithProgramHandler(
		teaHandler,
		termenv.ANSI256,
	)
}

func getEntries(url string) []*gofeed.Item {
	fp := gofeed.NewParser()
	// feed, _ := fp.ParseURL("https://github.com/Raven-bit/dotfiles/commits/primus.atom")
	feed, _ := fp.ParseURL(url)

	return feed.Items
}
