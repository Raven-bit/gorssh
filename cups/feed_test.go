package feed

// feed.go shouldn't be responsible for fetching http content at all.
// tea.Cmd should handle i/o. We really have our data structure living
// inside of our logic here ( which isn't great ), and we need to
// separate it .

// 1) we need to be fetching url content, and this should be pretty
// isolated.
// 2) We need to build a bunch of feed objects
// 3) these objects should be fed to views (cups of tea) in various ways.

// Step 1: Remove network I/O from the individual bubbletea programs
// step 2: write test for individual focus feed viewer program
// Step 3: ???
// Step 4: Profit!

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

func generateModel(feedstr string, t *testing.T) feedmodel {
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(RssFeed)
	if err != nil {
		t.Fatalf("Something went very wrong! %v\n", err)
	}
	if feed.Title != "Phoronix" {
		t.Fatalf("Not parsed right? Got '%v', expected 'Phoronix'.'", feed.Title)
	}
	model := NewFeedModel(feed.Items)
	return model
}

func TestDefaultCursor(t *testing.T) {
	model := generateModel(RssFeed, t)
	str := strings.Split(model.View(), "\n")

	if !strings.HasPrefix(str[2], ">") {
		t.Fatalf("First item in list is not selected? View is:\n%s", str)
	}
}

func TestCursorMoveDown(t *testing.T) {
	model := generateModel(RssFeed, t)
	modelcp, _ := model.Update(tea.KeyMsg(tea.Key{Type: tea.KeyRunes, Runes: []rune{'j'}}))
	model = modelcp.(feedmodel)
	str := strings.Split(model.View(), "\n")

	if !strings.HasPrefix(str[3], ">") {
		t.Fatalf("First item in list is not selected? View is:\n%s", str)
	}
}

func TestHeadLine(t *testing.T) {
	result := strings.Split(generateModel(RssFeed, t).View(), "\n")[0]
	// t.Logf("%v\n", str)
	expected := "What would you like to read?"

	if result != expected {
		t.Fatalf("First line strange, got: '%s'\nWanted: '%s'\n", result, expected)
	}
}

var RssFeed string = `
<rss version="2.0">
  <channel>
    <title>Phoronix</title>
    <link>https://www.phoronix.com/</link>
    <description>Linux Hardware Reviews & News</description>
    <language>en-us</language>
    <item>
      <title>NixOS 22.05 Released With New Graphical Installer</title>
      <link>
      https://www.phoronix.com/scan.php?page=news_item&px=NixOS-22.05-Released
      </link>
      <guid>
      https://www.phoronix.com/scan.php?page=news_item&px=NixOS-22.05-Released
      </guid>
      <description>
      NixOS as the Linux distribution built around the unique Nix package manager is out with its first release of the year...
      </description>
      <pubDate>Mon, 30 May 2022 17:05:49 -0400</pubDate>
    </item>
    <item>
      <title>
      OpenJPH 0.9 Released For Further Speeding Up Open-Source High-Throughput JPEG 2000
      </title>
      <link>
      https://www.phoronix.com/scan.php?page=news_item&px=OpenJPH-0.9
      </link>
      <guid>
      https://www.phoronix.com/scan.php?page=news_item&px=OpenJPH-0.9
      </guid>
      <description>
      While JPEG XL is regarded as the next-generation JPEG standard and JPEG 2000 never quite took off to supersede the original JPEG standard, there are open-source projects continuing to work on this image compression standard. OpenJPH 0.9 was released last week as the open-source high-throughput JPEG 2000 implementation and with this new version comes even more performance gains...
      </description>
      <pubDate>Mon, 30 May 2022 07:36:42 -0400</pubDate>
    </item>
    <item>
      <title>
      Raspberry Pi Sense HAT Joystick Driver Lands In Linux 5.19
      </title>
      <link>
      https://www.phoronix.com/scan.php?page=news_item&px=Raspberry-Pi-HAT-Joystick-Lands
      </link>
      <guid>
      https://www.phoronix.com/scan.php?page=news_item&px=Raspberry-Pi-HAT-Joystick-Lands
      </guid>
      <description>
      This weekend Linus Torvalds landed the Raspberry Pi Sense HT Joystick driver into the Linux 5.19 kernel as part of the input subsystem updates...
      </description>
      <pubDate>Mon, 30 May 2022 07:02:41 -0400</pubDate>
    </item>
  </channel>
</rss>
`
