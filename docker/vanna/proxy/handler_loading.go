package proxy

import (
	"fmt"
	"log"
	"net/http"
)

const loadingTemplate = `
<html>
<head>
    <title>Loading</title>
    <script type="text/javascript">
        function checkServerReady() {
            fetch(window.location.href, { method: 'GET' })
                .then(response => {
                    if (response.ok) {
                        // server is ready
                        window.location.reload();
                    } else {
                        // try to get the Retry-After header or default to 3000 ms
                        const retryAfter = parseInt(response.headers.get('Retry-After'), 10) * 1000 || 3000;
                        setTimeout(checkServerReady, retryAfter);
                    }
                })
                .catch(error => {
                    console.error('Error checking server readiness:', error);
                    // default to 3000 ms in case of other errors
                    setTimeout(checkServerReady, 3000);
                });
        }

        document.addEventListener("DOMContentLoaded", function() {
            setTimeout(checkServerReady, 3000);
        });
    </script>
</head>
<body>
    <h1>Loading! please wait...</h1>
</body>
</html>
`

func loadingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoadingHandler: %s\n", r.URL.Path)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, err := fmt.Fprint(w, loadingTemplate)
	if err != nil {
		log.Printf("Error writing loading template: %v\n", err)
	}
}

func serviceUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ServiceUnavailableHandler: %s\n", r.URL.Path)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Retry-After", "3")
	w.WriteHeader(http.StatusServiceUnavailable)
}
