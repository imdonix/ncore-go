package main
import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imdonix/ncore-go/pkg/ncore"
)

func main() {
	// ... (rest of the file remains same, but using ncore. instead of ncore.)
	username := os.Getenv("NCORE_USERNAME")
	password := os.Getenv("NCORE_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set NCORE_USERNAME and NCORE_PASSWORD environment variables")
	}

	// 1. Initialize the client
	client, err := ncore.NewClient(15*time.Second, nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 2. Login
	fmt.Printf("Logging in as %s...\n", username)
	_, err = client.Login(username, password, "")
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Login successful!")

	// 3. Search for a torrent
	searchTerm := "Forrest Gump"
	fmt.Printf("Searching for '%s'...\n", searchTerm)
	result, err := client.Search(searchTerm, ncore.TypeHDHun, ncore.WhereName, ncore.SortSeeders, ncore.SeqDesc, 1)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}


	if len(result.Torrents) == 0 {
		fmt.Println("No torrents found.")
	} else {
		fmt.Printf("Found %d torrents. Top result:\n", len(result.Torrents))
		top := result.Torrents[0]
		fmt.Printf("  Title: %s\n", top.Title)
		fmt.Printf("  Size:  %s\n", top.Size)
		fmt.Printf("  Seed:  %d | Leech: %d\n", top.Seeders, top.Leechers)
		fmt.Printf("  URL:   %s\n", top.URL)

		// 4. Example: Download the first result (commented out to avoid cluttering disk)

			fmt.Println("Downloading torrent file...")
			path, err := client.Download(top, ".", false)
			if err != nil {
				log.Printf("Download failed: %v", err)
			} else {
				fmt.Printf("Downloaded to: %s\n", path)
			}

	}

	// 5. Get activity (Hit & Run)
	fmt.Println("\nChecking Hit & Run activity...")
	activity, err := client.GetByActivity()
	if err != nil {
		log.Printf("Failed to get activity: %v", err)
	} else {
		fmt.Printf("You have %d active torrents in Hit & Run.\n", len(activity))
	}

	// 6. Logout
	client.Logout()
	fmt.Println("\nLogged out.")
}
