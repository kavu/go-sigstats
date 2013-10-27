# SigStats

Package sigstats is a tiny handler function for printing some statistics (mostly memory) for running application with just a tiny call of **SIGINFO** - **^T** shortcut.

You can view documentation on [godoc.org](http://godoc.org/github.com/kavu/go-sigstats "go-sigstats documentation").

## Example

```go

package main

import "github.com/kavu/go-sigstats"

func main() {
  // Let's enable all the Stats
  sigstats.EnableSigStats(&sigstats.SigStatsSettings{
    NumCPU:       true,
    NumCGOCalls:  true,
    NumGoroutine: true,
    MemStats:     true,
    GoVersion:    true,
    GCStats:      true,
  })

  select {}
}
```

### HTTP Server Requests Tracker example

I am using [gorilla/mux](https://github.com/gorilla/mux) but you can try default `ServeMux` from `net/http`.

```go
package main

import (
	"github.com/gorilla/mux"
	"github.com/kavu/go-sigstats"
	"io"
	"log"
	"net/http"
	"time"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 3)
	io.WriteString(w, "hello, world!\n")
}

func main() {
	tracker := sigstats.InitSigStatsHTTPRequstTracker()

	mx := mux.NewRouter()
	mx.HandleFunc("/", HelloServer)

	log.Fatal(http.ListenAndServe(":8080", sigstats.SigStatsTrackerMiddleware(mx, tracker)))
}

```

Run this example and in parallel terminal window start some stress loader tool like [wrk](https://github.com/wg/wrk) and press **^T** several times.

## Thanks

Inspired by [Dustin Sallings](https://github.com/dustin) [post](http://dustin.github.io/2013/07/04/siginfo.html) about SIGINFO. 
