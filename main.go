package main

import (
	"fmt"
	"log"
	"os"

	// "io"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"

	feed "github.com/mmcdole/gofeed"
)

type feeditems struct {
	items    []*feed.Item
	cursor   int
	viewing  bool
	selected *feed.Item
}

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

func main() {
	// fmt.Printf("%v\n", initialFeedItems())
	p := tea.NewProgram(initialFeedItems())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (f feeditems) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			if f.viewing {
				f.viewing = false
			} else {
				return f, tea.Quit
			}

			// The "up" and "k" keys move the cursor up
		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
			}

			// The "down" and "j" keys move the cursor down
		case "down", "j":
			if f.cursor < len(f.items)-1 {
				f.cursor++
			}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			// case "enter", " ":
			//     _, ok := f.selected[f.cursor]
			//     if ok {
			//         delete(f.selected, m.cursor)
			//     } else {
			//         f.selected[m.cursor] = struct{}{}
			//     }
			// }
		case "enter", " ":
			f.viewing = true
			f.selected = f.items[f.cursor]
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return f, nil
}

func (f feeditems) View() string {
	if f.viewing {
		return f.viewItem()
	}

	return f.viewList()
}

func (f feeditems) viewList() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range f.items {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if f.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		// checked := " " // not selected
		// if _, ok := m.selected[i]; ok {
		//     checked = "x" // selected!
		// }

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice.Title)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (f feeditems) viewItem() string {
	var s string
	s += fmt.Sprintf("Title: %s\n\n", f.selected.Title)
	s += fmt.Sprintf("By: ")
	numAuthors := len(f.selected.Authors)
	if numAuthors != 0 {
		if numAuthors > 1 {
			for i, author := range f.selected.Authors {
				if numAuthors == i {
					s += ", and "
				}

				s += author.Name

				if numAuthors > i {
					s += ", "
				}
			}
		} else if numAuthors == 1 {
			s += fmt.Sprintf("%s\n\n", f.selected.Authors[0].Name)
		}
	}

	s += fmt.Sprintf("Associated URL: %s\n\n", f.selected.Link)

	converter := md.NewConverter("", true, nil)
	cnt, err := converter.ConvertString(f.selected.Description)
	if err != nil {
		log.Fatal(err)
	}
	cnt, err = glamour.Render(cnt, "dark")
	if err != nil {
		log.Fatal(err)
	}
	s += fmt.Sprintf("\n%s", cnt)
	return s
}

func initialFeedItems() feeditems {
	return feeditems{
		items:    getEntries("https://www.phoronix.com/rss.php"),
		selected: &feed.Item{},
	}
}

func getEntries(url string) []*feed.Item {
	fp := feed.NewParser()
	// feed, _ := fp.ParseURL("https://github.com/Raven-bit/dotfiles/commits/primus.atom")
	feed, _ := fp.ParseURL(url)

	return feed.Items
}

func (f feeditems) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
