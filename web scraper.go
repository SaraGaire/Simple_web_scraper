package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Article represents a scraped article
type Article struct {
	Title string
	Link  string
}

// scrapeHackerNews scrapes the front page of Hacker News
func scrapeHackerNews() ([]Article, error) {
	// Make HTTP GET request
	resp, err := http.Get("https://news.ycombinator.com")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var articles []Article

	// Find and extract article titles and links
	doc.Find(".titleline > a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, exists := s.Attr("href")
		
		if exists && title != "" {
			// Handle relative URLs
			if strings.HasPrefix(link, "item?id=") {
				link = "https://news.ycombinator.com/" + link
			}
			
			articles = append(articles, Article{
				Title: title,
				Link:  link,
			})
		}
	})

	return articles, nil
}

// scrapeQuotes scrapes quotes from a demo quotes website
func scrapeQuotes() ([]map[string]string, error) {
	resp, err := http.Get("http://quotes.toscrape.com")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quotes page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var quotes []map[string]string

	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		quote := strings.TrimSpace(s.Find(".text").Text())
		author := strings.TrimSpace(s.Find(".author").Text())
		
		if quote != "" && author != "" {
			quotes = append(quotes, map[string]string{
				"quote":  quote,
				"author": author,
			})
		}
	})

	return quotes, nil
}

// Generic scraper function that can be customized
func genericScraper(url string, selector string) ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var results []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})

	return results, nil
}

func main() {
	fmt.Println("🕷️  Go Web Scraper Demo")
	fmt.Println("=" * 40)

	// Example 1: Scrape Hacker News
	fmt.Println("\n📰 Scraping Hacker News headlines...")
	articles, err := scrapeHackerNews()
	if err != nil {
		log.Printf("Error scraping Hacker News: %v", err)
	} else {
		fmt.Printf("Found %d articles:\n\n", len(articles))
		for i, article := range articles[:5] { // Show only first 5
			fmt.Printf("%d. %s\n   🔗 %s\n\n", i+1, article.Title, article.Link)
		}
	}

	// Example 2: Scrape quotes
	fmt.Println("\n💭 Scraping quotes...")
	quotes, err := scrapeQuotes()
	if err != nil {
		log.Printf("Error scraping quotes: %v", err)
	} else {
		fmt.Printf("Found %d quotes:\n\n", len(quotes))
		for i, quote := range quotes[:3] { // Show only first 3
			fmt.Printf("%d. \"%s\"\n   - %s\n\n", i+1, quote["quote"], quote["author"])
		}
	}

	// Example 3: Generic scraper
	fmt.Println("\n🔧 Using generic scraper for GitHub trending...")
	trending, err := genericScraper("https://github.com/trending", "h2.h3 a")
	if err != nil {
		log.Printf("Error with generic scraper: %v", err)
	} else {
		fmt.Printf("Found %d trending repositories:\n\n", len(trending))
		for i, repo := range trending[:5] { // Show only first 5
			fmt.Printf("%d. %s\n", i+1, repo)
		}
	}

	fmt.Println("\n✅ Scraping completed!")
}
