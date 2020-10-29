package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"web-app/pubsub/publisher"
	"web-app/pubsub/subscriber"

	"web-app/utils"
)

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

const maxMessages = 10

func main() {
    // Load .env
    utils.Load()

    // Initialise pubsub
	publisher.Init()

    // Declare http routing
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/pubsub/publish", publishHandler)
	http.HandleFunc("/pubsub/push", pushHandler)
	http.HandleFunc("/pubsub/subscribe", subscribeHandler)

    // Listen and Serve
    port := utils.Getenv("PORT", "8080")
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", "Subscribe request", r.RequestURI)

	projectID := utils.MustGetenv("GOOGLE_CLOUD_PROJECT")
	subID := utils.MustGetenv("PUBSUB_SUBSCRIPTION_ID")

	subscriber.PullMsgsConcurrenyControl(log.Writer(), projectID, subID)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", "Push request", r.RequestURI)

	// Verify the token.
	if r.URL.Query().Get("token") != utils.MustGetenv("PUBSUB_VERIFICATION_TOKEN") {
		http.Error(w, "Bad token", http.StatusBadRequest)
	}

	msg := publisher.CreatePushRequest()
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	messages = append(messages, string(msg.Message.Data))
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", "List request", r.RequestURI)

	messagesMu.Lock()
	defer messagesMu.Unlock()

	if err := tmpl.Execute(w, messages); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", "Publish request", r.RequestURI)

	if _, err := publisher.SendMessage(r.FormValue("payload")); err != nil {
		http.Error(w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}
	fmt.Fprint(w, "Message published.")
}

var tmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>Pub/Sub</title>
  </head>
  <body>
    <div>
      <p>Last ten messages received by this instance:</p>
      <ul>
      {{ range . }}
          <li>{{ . }}</li>
      {{ end }}
      </ul>
    </div>
    <form method="post" action="/pubsub/publish">
      <textarea name="payload" placeholder="Enter message here"></textarea>
      <input type="submit">
    </form>
    <p>Note: if the application is running across multiple instances, each
      instance will have its own list of messages.</p>
  </body>
</html>`))
