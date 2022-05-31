package feed

import (
	"fmt"
	"log"

	// "io"

	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"

	feed "github.com/mmcdole/gofeed"
)

type feedmodel struct {
	items    []*feed.Item
	cursor   int
	viewing  bool
	selected *feed.Item
}

func (f feedmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (f feedmodel) View() string {
	if f.viewing {
		return f.viewItem()
	}
	return f.viewList()
}

func (f feedmodel) viewList() string {
	// The header
	s := "What would you like to read?\n\n"

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

func (f feedmodel) viewItem() string {
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

func NewFeedModel(items []*feed.Item) feedmodel {
	return feedmodel{
		items:    items,
		selected: &feed.Item{},
	}
}

func (f feedmodel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
