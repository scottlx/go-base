package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Fetch() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error, fetch: %v\n", err)
			os.Exit(1)
		}
		n, err := io.Copy(os.Stdout, resp.Body)
		fmt.Println("write", n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading body: %v\n", err)
			os.Exit(1)
		}
		resp.Body.Close()
	}
}
