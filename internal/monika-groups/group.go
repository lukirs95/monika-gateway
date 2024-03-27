package monikagroups

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	monikaauth "github.com/lukirs95/monika-gateway/internal/monika-auth"
	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db/dbrepo"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

type GroupHandler struct {
	db monikadb.GroupDatabaseRepo
}

func NewGroupHandler(db monikadb.GroupDatabaseRepo) *GroupHandler {
	return &GroupHandler{
		db: db,
	}
}

var (
	regGroups  = regexp.MustCompile(`^/.+/groups$`)
	regGroupId = regexp.MustCompile(`^/.+/groups/(\d+)$`)
)

func (h *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := monikaauth.GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if regGroups.MatchString(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetGroups(w, r)
			return
		case http.MethodPost:
			if session.Role != sdk.UserRole_ADMIN {
				http.Error(w, "userlevel too low", http.StatusUnauthorized)
				return
			}
			h.handleNewGroup(w, r)
			return
		default:
			http.Error(w, "Method not implemented", http.StatusMethodNotAllowed)
			return
		}
	}

	if regGroupId.MatchString(r.URL.Path) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetGroup(w, r)
			return
		case http.MethodPatch:
			if session.Role != sdk.UserRole_ADMIN {
				http.Error(w, "userlevel too low", http.StatusUnauthorized)
				return
			}
			h.handleUpdateGroup(w, r)
			return
		case http.MethodDelete:
			if session.Role != sdk.UserRole_ADMIN {
				http.Error(w, "userlevel too low", http.StatusUnauthorized)
				return
			}
			h.handleDeleteGroup(w, r)
			return
		default:
			http.Error(w, "Method not implemented", http.StatusMethodNotAllowed)
			return
		}
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *GroupHandler) handleNewGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup sdk.Group
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.CreateGroup(&newGroup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&newGroup); err != nil {
		log.Print(err)
		return
	}
}

func (h *GroupHandler) handleGetGroups(w http.ResponseWriter, _ *http.Request) {
	groups, err := h.db.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&groups); err != nil {
		log.Print(err)
		return
	}
}

func (h *GroupHandler) handleGetGroup(w http.ResponseWriter, r *http.Request) {
	parsedGroupId := regGroupId.FindStringSubmatch(r.URL.Path)[1]
	groupId, err := strconv.ParseInt(parsedGroupId, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse group id", http.StatusBadRequest)
		return
	}

	group, err := h.db.GetGroup(groupId)
	if err != nil {
		http.Error(w, "no group found", http.StatusBadRequest)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&group); err != nil {
		log.Print(err)
		return
	}
}

func (h *GroupHandler) handleUpdateGroup(w http.ResponseWriter, r *http.Request) {
	parsedGroupId := regGroupId.FindStringSubmatch(r.URL.Path)[1]
	groupId, err := strconv.ParseInt(parsedGroupId, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse group id", http.StatusBadRequest)
		return
	}

	var group sdk.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateGroupname(groupId, group.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *GroupHandler) handleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	parsedGroupId := regGroupId.FindStringSubmatch(r.URL.Path)[1]
	groupId, err := strconv.ParseInt(parsedGroupId, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse group id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteGroup(groupId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
