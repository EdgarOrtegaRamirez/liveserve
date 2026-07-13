package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// injectLiveReload reads an HTML file, injects the live reload script,
// and writes it to the response writer.
func injectLiveReload(w http.ResponseWriter, r *http.Request, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	content := string(data)

	// Inject the live reload script before </body> or at the end
	reloadScript := fmt.Sprintf(`<script>
(function() {
	var ws = new WebSocket("ws://%s/__livereload");
	ws.onmessage = function(e) {
		if (e.data === "reload") {
			location.reload();
		}
	};
})();
</script>
</body>`, r.Host)

	if strings.Contains(content, "</body>") {
		content = strings.Replace(content, "</body>", reloadScript, 1)
	} else {
		content += reloadScript
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, content)
}