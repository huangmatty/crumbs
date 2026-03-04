package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<html>
		<body>
			<h1>Crumbs Admin Portal</h1>
			<p>Crumbs has been visited %d times.</p>
		</body>
	</html>
	`, cfg.fsHits.Load())
}
