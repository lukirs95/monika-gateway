package monikadevice

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	monikadriver "github.com/lukirs95/monika-gateway/internal/monika-driver"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type MonikaDeviceHandler struct {
	driverController *monikadriver.DriverController
}

func NewMonikaDeviceHandler(driverController *monikadriver.DriverController) *MonikaDeviceHandler {
	return &MonikaDeviceHandler{
		driverController: driverController,
	}
}

var (
	regDevices = regexp.MustCompile(`^/.+/devices$`)
	regDevice  = regexp.MustCompile(`^/.+/devices/`)
)

func (handler *MonikaDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if regDevices.MatchString(r.URL.Path) {
		handler.handleGetAllDevices(w)
		return
	}
	if regDevice.MatchString(r.URL.Path) {
		handler.proxyToDriver(w, r)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func (handler *MonikaDeviceHandler) handleGetAllDevices(w http.ResponseWriter) {
	devices := make([]types.Device, 0)
	for _, location := range handler.driverController.GetLocations() {
		res, err := http.Get(fmt.Sprintf("http://%s", location))
		if err != nil {
			log.Print(err)
			continue
		}

		decoder := json.NewDecoder(res.Body)

		nextDevices, err := types.DevicesFromJSON(decoder)
		if err != nil {
			log.Print(err)
			continue
		}
		devices = append(devices, nextDevices...)
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&devices); err != nil {
		log.Print(err)
	}
}

func (handler *MonikaDeviceHandler) proxyToDriver(w http.ResponseWriter, r *http.Request) {
	_, after, found := strings.Cut(r.URL.Path, "/devices/")
	if !found {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	composite := strings.SplitN(after, "/", 2)
	deviceType := types.DeviceType(composite[0])
	path := ""
	if len(composite) > 1 {
		path = composite[1]
	}

	driverLocation := handler.driverController.GetLocation(deviceType)
	if driverLocation == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	requestPath := fmt.Sprintf("http://%s/%s", driverLocation, path)
	requestMethod := r.Method

	if requestMethod == http.MethodGet {
		res, err := http.Get(requestPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", res.Header.Get("Content-Type"))
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
		return
	}

	if requestMethod == http.MethodPost {
		res, err := http.Post(requestPath, r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", res.Header.Get("Content-Type"))
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
		return
	}

	http.Error(w, "not allowed", http.StatusMethodNotAllowed)
}
