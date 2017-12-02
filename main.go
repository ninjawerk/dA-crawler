package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
)

type Artwork struct {
	gorm.Model
	Title     string
	Artist    string
	Url       string
	ImageUrl  string
	FavCount  int
	ArtistUrl string
}

type Link struct {
	gorm.Model
	Url string
}

var (
	// Protect access to dup
	mu sync.Mutex
	// Duplicates table
	dup = map[string]bool{}
	db  *gorm.DB
	// Command-line flags
	seed = flag.String("seed", "https://www.deviantart.com/", "seed URL")
)

func main() {
	//init the db
	var dbErr error
	db, dbErr = gorm.Open("postgres", "host=localhost user=postgres dbname=dA_Data sslmode=disable password=")
	if dbErr != nil {
		fmt.Println(dbErr)
	}
	db.AutoMigrate(&Artwork{})
	db.AutoMigrate(&Link{})
	defer db.Close()

	flag.Parse()

	// Parse the provided seed
	_, err := url.Parse(*seed)
	if err != nil {
		log.Fatal(err)
	}

	// Create the muxer
	mux := fetchbot.NewMux()

	// Handle all errors the same
	mux.HandleErrors(fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
	}))

	// Handle GET requests for html responses, to parse the body and enqueue all links as HEAD
	// requests.
	mux.Response().Method("GET").ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			// Process the body to find the links
			doc, err := goquery.NewDocumentFromResponse(res)

			if err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				return
			}

			favstr := doc.Find("div.who-faved-modal").First()
			favtext := favstr.Parent().First().Text()

			favtextReplaced := strings.Replace(favtext, ",", "", -1)
			favCount, _ := strconv.Atoi(strings.Split(favtextReplaced, " (who?)")[0])

			imgContainer := doc.Find("div.dev-view-deviation").First()
			img := imgContainer.Find("img.dev-content-full").First()

			titleContainer := doc.Find("div.dev-title-container ").First()
			artUrl, _ := titleContainer.Find("a.title").First().Attr("href")
			title := titleContainer.Find("a.title").First().Text()
			authorName := titleContainer.Find("a.username").First().Text()
			authorLink, _ := titleContainer.Find("a.username").First().Attr("href")
			imageSrc, _ := img.Attr("src")
			//just double check if its in the db again, sometimes dA directs to Url from aliases
			var count int
			db.Table("artworks").Where("Url = ?", artUrl).Count(&count)
			if count <= 0 {
				artwork := Artwork{Title: title, Url: artUrl, Artist: authorName, ArtistUrl: authorLink, ImageUrl: imageSrc, FavCount: favCount }
				db.Create(&artwork)
			}

			if err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				return
			}

			// Enqueue all links as HEAD requests
			enqueueLinks(ctx, doc)
		}))

	// Create the Fetcher, handle the logging first, then dispatch to the Muxer
	h := logHandler(mux)

	f := fetchbot.New(h)

	// Start processing
	q := f.Start()

	// Enqueue the seed, which is the first entry in the dup map
	dup[*seed] = true
	_, err = q.SendStringGet(*seed)
	if err != nil {
		fmt.Printf("[ERR] GET %s - %s\n", *seed, err)
	}
	q.Block()
}

// logHandler prints the fetch information and dispatches the call to the wrapped Handler.
func logHandler(wrapped fetchbot.Handler) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		if err == nil {
			fmt.Printf("[%d] %s %s - %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL(), res.Header.Get("Content-Type"))
		}
		wrapped.Handle(ctx, res, err)
	})
}

func enqueueLinks(ctx *fetchbot.Context, doc *goquery.Document) {
	mu.Lock()
	fmt.Println("Processing for links... " + ctx.Cmd.URL().String())
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		u, err := ctx.Cmd.URL().Parse(val)

		// get out if its not an artwork
		if !strings.Contains(val, "/art/") || !strings.Contains(u.Host, "deviantart.com") {
			return
		}

		if err != nil {
			fmt.Printf("error: resolve URL %s - %s\n", val, err)
			return
		}
		var count int
		db.Table("links").Where("Url = ?", u.String()).Count(&count)
		if count <= 0 {
			if _, err := ctx.Q.SendStringGet(u.String()); err != nil {
				fmt.Printf("error: enqueue head %s - %s\n", u, err)
			} else {
				nlink := Link{Url: u.String()}
				db.Create(&nlink)
			}
		}

		doc = nil
	})
	mu.Unlock()
}
