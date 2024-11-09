package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// GetClientIP retrieves the client's IP address from the request
func GetClientIP(r *http.Request) string {
	// Check for the "X-Forwarded-For" header if the request is behind a proxy
	// (e.g., if running in a cloud environment or using a reverse proxy like Nginx)
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		// If the header is not available, use the remote address
		clientIP = r.RemoteAddr
	}
	// Remove port from IP if it exists (e.g., "192.168.1.1:8080" -> "192.168.1.1")
	clientIP = strings.Split(clientIP, ":")[0]
	return clientIP
}

// chromedpTask runs a ChromeDP task to open a specified URL and keeps it open for a duration
func chromedpTask(w http.ResponseWriter, r *http.Request) {
	// Retrieve the client's IP address
	clientIP := GetClientIP(r)

	// Generate a dynamic URL or content based on the client's IP address
	dynamicURL := fmt.Sprintf("https://example.com/?ip=%s", clientIP)

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

	// Keep the page open for a specific duration (e.g., 2 minutes)
	time.Sleep(2 * time.Minute) // Keeps the browser open for 2 minutes (adjust as needed)

	// Send a response to the client with dynamic content
	fmt.Fprintf(w, "Client IP: %s\n", clientIP)
	fmt.Fprintf(w, "Dynamic URL based on your IP: %s\n", dynamicURL)
}

func main() {
	// Set up HTTP route
	http.HandleFunc("/", chromedpTask)

	// Start the HTTP server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
