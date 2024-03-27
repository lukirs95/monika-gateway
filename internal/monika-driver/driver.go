package monikadriver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type DriverController struct {
	sync.Mutex
	drivers []types.Driver
}

func NewDriverController() *DriverController {
	return &DriverController{
		drivers: make([]types.Driver, 0),
	}
}

var (
	regDriverConnect   = regexp.MustCompile(`^/driver/connect$`)
	regDriverDisonnect = regexp.MustCompile(`^/driver/disconnect$`)
)

func (controller *DriverController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if regDriverConnect.MatchString(r.URL.Path) {
		controller.handleConnect(w, r)
		return
	}
	if regDriverDisonnect.MatchString(r.URL.Path) {
		controller.handleDisconnect(w, r)
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (controller *DriverController) handleConnect(w http.ResponseWriter, r *http.Request) {
	var newDriver types.Driver
	if err := json.NewDecoder(r.Body).Decode(&newDriver); err != nil {
		http.Error(w, "could not read json", http.StatusBadRequest)
		return
	}

	if newDriver.Location == "" {
		newDriver.Location = strings.Split(r.RemoteAddr, ":")[0]
	}

	controller.Lock()
	controller.drivers = append(controller.drivers, newDriver)
	controller.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func (controller *DriverController) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	var newDriver types.Driver
	if err := json.NewDecoder(r.Body).Decode(&newDriver); err != nil {
		http.Error(w, "could not read json", http.StatusBadRequest)
		return
	}

	for index, driver := range controller.drivers {
		if driver.DeviceType == newDriver.DeviceType {
			controller.Lock()
			controller.drivers = append(controller.drivers[:index], controller.drivers[index:]...)
			controller.Unlock()
		}
	}
}

func (controller *DriverController) GetLocation(deviceType types.DeviceType) string {
	controller.Lock()
	defer controller.Unlock()
	for _, driver := range controller.drivers {
		if driver.DeviceType == deviceType {
			return fmt.Sprintf("%s:%d", driver.Location, driver.Port)
		}
	}
	return ""
}

func (controller *DriverController) GetLocations() []string {
	locations := make([]string, 0)
	for _, driver := range controller.drivers {
		locations = append(locations, fmt.Sprintf("%s:%d", driver.Location, driver.Port))
	}
	return locations
}
