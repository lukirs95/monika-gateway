package monikaauth

import (
	"encoding/json"
	"log"
	"net/http"

	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (auth *MonikaAuth) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	session, err := GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if session.Role != sdk.UserRole_ADMIN {
		http.Error(w, "userlevel too low", http.StatusUnauthorized)
		return
	}

	var askingUser sdk.User
	if err := json.NewDecoder(r.Body).Decode(&askingUser); err != nil {
		log.Println(err)
		http.Error(w, "could not parse as json", http.StatusBadRequest)
		return
	}

	if err := auth.CreateUser(&askingUser); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (auth *MonikaAuth) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session, err := GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if session.Role != sdk.UserRole_ADMIN {
		http.Error(w, "userlevel too low", http.StatusUnauthorized)
		return
	}

	users, err := auth.GetAllUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, "could not find user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println(err)
		http.Error(w, "could not encode user", http.StatusInternalServerError)
		return
	}
}

func (auth *MonikaAuth) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	session, err := GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := auth.GetUserById(session.UserId)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not find user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		http.Error(w, "could not encode user", http.StatusInternalServerError)
		return
	}
}

func (auth *MonikaAuth) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var askingUser sdk.User
	if err := json.NewDecoder(r.Body).Decode(&askingUser); err != nil {
		log.Println(err)
		http.Error(w, "could not parse as json", http.StatusBadRequest)
		return
	}

	session, err := GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if askingUser.Password != "" {
		if askingUser.UserId == 0 {
			askingUser.UserId = session.UserId
		}

		if session.Role == sdk.UserRole_ADMIN || session.UserId == askingUser.UserId {
			if err := auth.UpdateUserPassword(&askingUser); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "not allowed to update user password", http.StatusUnauthorized)
			return
		}
	}

	if askingUser.Role != "" {
		if session.Role == sdk.UserRole_ADMIN {
			if err := auth.UpdateUserRole(askingUser); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "not allowed to update user role", http.StatusUnauthorized)
			return
		}
	}
}

func (auth *MonikaAuth) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var askingUser sdk.User
	if err := json.NewDecoder(r.Body).Decode(&askingUser); err != nil {
		log.Println(err)
		http.Error(w, "could not parse as json", http.StatusBadRequest)
		return
	}

	session, err := GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// admins can delete any user, normal users can only delete themselves
	if session.Role == sdk.UserRole_ADMIN || session.UserId == askingUser.UserId {
		if err := auth.DeleteUser(askingUser); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	http.Error(w, "userlevel too low", http.StatusUnauthorized)
}
