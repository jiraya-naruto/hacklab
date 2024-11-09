package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

// chromedpTask runs a ChromeDP task to open a specified URL and keeps it open for a duration
func chromedpTask(w http.ResponseWriter, r *http.Request) {
	// URL to navigate to
	url := "https://jiraya-naruto.github.io/jiraya/" // Replace with your webpage URL

	// Set up Chrome options
	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),                      // Ensure headless mode is off
		chromedp.Flag("start-fullscreen", true),               // Full-screen mode
		chromedp.Flag("disable-infobars", true),               // Disable "Chrome is being controlled" message
		chromedp.Flag("disable-features", "TranslateUI"),      // Disable translate UI
		chromedp.Flag("kiosk", true),                          // Kiosk mode (borderless fullscreen)
		chromedp.Flag("disable-ui-for-tests", true),           // Disable UI for tests
		chromedp.Flag("overscroll-history-navigation", false), // Disable scroll navigation
		chromedp.Flag("no-default-browser-check", true),       // Disable default browser check
		chromedp.Flag("disable-pinch", true),                  // Disable pinch zoom
	)

	// Set up context with the specified Chrome options
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// Create a new ChromeDP context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Run the ChromeDP task
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		log.Println("Failed to load page:", err)
		return
	}

	// Keep the page open for a specific duration
	time.Sleep(864000 * time.Second) // Keeps the browser open for 2 minutes (adjust as needed)

	// Send a response to the client after the duration
	fmt.Fprintln(w, "Page was loaded and kept open for 2 hours")
}

func main() {
	// Set up HTTP route
	http.HandleFunc("/", chromedpTask)

	// Start the HTTP server
	fmt.Println("Starting server on :3000")
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
