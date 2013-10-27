// Copyright (C) 2013 Max Riveiro
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sigstats

// Just a lockable map with current running requests
type httpServerTracker struct {
	sync.RWMutex
	requests map[*http.Request]time.Time
}

// Add request to Tracker
func (self *httpServerTracker) add(r *http.Request, t time.Time) {
	self.Lock()
	defer self.Unlock()
	self.requests[r] = t
}

// Remove completed request from Tracker
func (self *httpServerTracker) remove(r *http.Request) {
	self.Lock()
	defer self.Unlock()
	delete(self.requests, r)
}

// Print a list of the current running http.Request
func (self *httpServerTracker) report() {
	self.RLock()
	defer self.RUnlock()
	if len(self.requests) > 0 {
		fmt.Println("\nCurrent active HTTP Requests:")
		for r, t := range self.requests {
			fmt.Printf("\t%s %s %s for %v\n", r.RemoteAddr, r.Method, r.URL, time.Since(t))
		}
	} else {
		fmt.Println("\nThere is no active HTTP Requests.")
	}
	fmt.Println()
}

// We need to init Sigstats Handler and initilize new Tracker for actually tracking
func InitSigStatsHTTPRequstTracker() *httpServerTracker {
	tr := &httpServerTracker{
		requests: map[*http.Request]time.Time{},
	}

	InitInfoHandler(tr.report)

	return tr
}

// HTTP Middleware. Can be used with different muxers, such as https://github.com/gorilla/mux. You need to pass a pointer to httpServerTracker, where middleware will collect runtime data.
func SigStatsTrackerMiddleware(h http.Handler, tr *httpServerTracker) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		tr.add(r, time.Now())
		h.ServeHTTP(w, r)
		tr.remove(r)
	}
	return http.HandlerFunc(f)
}
