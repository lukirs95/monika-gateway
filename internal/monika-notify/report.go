package monikanotify

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type report struct {
	Type  ReportType `json:"type"`
	Value any        `json:"value"`
}

type ReportType string

const (
	ReportType_Error        ReportType = "ERROR"
	ReportType_ErrorRelease ReportType = "ERROR_RELEASE"
	ReportType_Update       ReportType = "UPDATE"
)

type Controller struct {
	subscriptions []*subscriber
	errors        map[int64]*types.PubError
	errorId       int64
}

func NewController() *Controller {
	controller := &Controller{
		subscriptions: make([]*subscriber, 0),
		errors:        make(map[int64]*types.PubError),
		errorId:       0,
	}

	return controller
}

var regError = regexp.MustCompile(`^/error$`)
var regErrorId = regexp.MustCompile(`^/error/\d+$`)

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/notify")
	log.Print(path)
	switch r.Method {
	case http.MethodGet:
		if path == "/ws" {
			log.Print("new websocket connection")
			subscriber := c.newSubscription()
			subscriber.upgrade(w, r)
			log.Print("delete websocket connection")
			c.deleteSubscription(subscriber)
			return
		}
	case http.MethodPost:
		if path == "/update" {
			c.notifyUpdate(w, r)
			return
		}
		if regError.MatchString(path) {
			c.notifyError(w, r)
			return
		}
	case http.MethodDelete:
		if regErrorId.MatchString(path) {
			c.deleteError(w, r)
			return
		}
	default:
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
}

func (c *Controller) newSubscription() *subscriber {
	subscriber := newSubscriber(c.errors)
	c.subscriptions = append(c.subscriptions, subscriber)
	return subscriber
}

func (c *Controller) deleteSubscription(subscriber *subscriber) {
	for i, subscription := range c.subscriptions {
		if subscription == subscriber {
			c.subscriptions = append(c.subscriptions[:i], c.subscriptions[i+1:]...)
		}
	}
}

func (c *Controller) notifyUpdate(w http.ResponseWriter, r *http.Request) {
	var device types.DeviceUpdate
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, subscription := range c.subscriptions {
		subscription.notifyUpdate(device)
	}
}

func (c *Controller) notifyError(w http.ResponseWriter, r *http.Request) {
	pubError := types.PubError{}
	if err := json.NewDecoder(r.Body).Decode(&pubError); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pubError.ErrorId = c.errorId
	c.errors[c.errorId] = &pubError
	for _, subscription := range c.subscriptions {
		subscription.notifyError(&pubError)
	}

	res := types.PubErrorResponse{ErrorId: c.errorId}

	w.WriteHeader(http.StatusCreated)
	header := w.Header()
	header.Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		log.Print(err)
	}
	c.errorId++
}

func (c *Controller) deleteError(w http.ResponseWriter, r *http.Request) {
	errorId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/notify/error/"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldError := c.errors[errorId]
	if oldError == nil {
		return
	}

	for _, subscription := range c.subscriptions {
		subscription.deleteError(oldError)
	}

	delete(c.errors, errorId)
}
