package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	monikaauth "github.com/lukirs95/monika-gateway/internal/monika-auth"
	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	monikadevice "github.com/lukirs95/monika-gateway/internal/monika-device"
	monikadriver "github.com/lukirs95/monika-gateway/internal/monika-driver"
	monikagroups "github.com/lukirs95/monika-gateway/internal/monika-groups"
	monikanotify "github.com/lukirs95/monika-gateway/internal/monika-notify"
)

var auth *monikaauth.MonikaAuth
var adminPassword = "admin"

func main() {
	defaultPassword, err := monikaauth.HashDefaultPassword(adminPassword)
	if err != nil {
		panic(err)
	}

	db, close := monikadb.NewSQLiteDatabase(":memory:", defaultPassword)
	defer close()
	db.Migrate()

	authParams := monikaauth.NewMonikaAuthParams(time.Hour, []byte("TEST"))
	auth = monikaauth.NewMonikaAuth(db, authParams)

	router := mux.NewRouter()
	router.HandleFunc("/api/authenticate", auth.HandleLogin)
	router.HandleFunc("/api/users/self", auth.Secure(auth.HandleGetUser)).Methods("GET")
	router.HandleFunc("/api/users", auth.Secure(auth.HandleGetAllUsers)).Methods("GET")
	router.HandleFunc("/api/users", auth.Secure(auth.HandleCreateUser)).Methods("POST")
	router.HandleFunc("/api/users", auth.Secure(auth.HandleUpdateUser)).Methods("PATCH")
	router.HandleFunc("/api/users", auth.Secure(auth.HandleDeleteUser)).Methods("DELETE")

	notifier := monikanotify.NewController()
	router.PathPrefix("/api/notify").Handler(notifier)

	memberHandler := monikagroups.NewMemberHandler(db)
	router.PathPrefix("/api/groups/members").HandlerFunc(auth.SecureHandler(memberHandler))
	groupHandler := monikagroups.NewGroupHandler(db)
	router.PathPrefix("/api/groups").HandlerFunc(auth.SecureHandler(groupHandler))
	driverHandler := monikadriver.NewDriverController()
	router.PathPrefix("/driver").Handler(driverHandler)
	deviceHandler := monikadevice.NewMonikaDeviceHandler(driverHandler)
	router.PathPrefix("/api/devices").Handler(deviceHandler)

	http.ListenAndServe(":8080", router)
}
