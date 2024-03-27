package monikagroups

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	monikaauth "github.com/lukirs95/monika-gateway/internal/monika-auth"
	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db/dbrepo"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

type MemberHandler struct {
	db monikadb.GroupDatabaseRepo
}

func NewMemberHandler(db monikadb.GroupDatabaseRepo) *MemberHandler {
	return &MemberHandler{
		db: db,
	}
}

var (
	regMembers  = regexp.MustCompile(`^/.+/members\?$`)
	regMemberId = regexp.MustCompile(`^/.+/members/(\d+)\?$`)
	regGroupsId = regexp.MustCompile(`^/.+/members\?(groupId=(\d+),?)$`)
)

func (h *MemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := monikaauth.GetSessionFrom(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	path := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)

	// /api/members
	if regMembers.MatchString(path) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetAllMembers(w, r)
			return
		case http.MethodPost:
			if session.Role != sdk.UserRole_ADMIN {
				http.Error(w, "userlevel too low", http.StatusUnauthorized)
				return
			}
			h.handleCreateMember(w, r)
			return
		default:
			http.Error(w, "Not implemented", http.StatusMethodNotAllowed)
			return
		}
	}

	// /api/members/{memberId}
	if regMemberId.MatchString(path) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetMember(w, r)
			return
		case http.MethodDelete:
			if session.Role != sdk.UserRole_ADMIN {
				http.Error(w, "userlevel too low", http.StatusUnauthorized)
				return
			}
			h.handleDeleteMember(w, r)
			return
		default:
			http.Error(w, "Not implemented", http.StatusMethodNotAllowed)
			return
		}
	}

	// /api/members?groupId={groupId}
	if regGroupsId.MatchString(path) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetMembersByGroup(w, r)
			return
		default:
			http.Error(w, "Not implemented", http.StatusMethodNotAllowed)
			return
		}
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *MemberHandler) handleCreateMember(w http.ResponseWriter, r *http.Request) {
	var member sdk.GroupMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.CreateMember(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&member); err != nil {
		log.Print(err)
		return
	}
}

func (h *MemberHandler) handleGetMember(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
	parsedMemberId := regMemberId.FindStringSubmatch(path)[1]
	memberId, err := strconv.ParseInt(parsedMemberId, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse module id", http.StatusBadRequest)
		return
	}

	member, err := h.db.GetMember(memberId)
	if err != nil {
		http.Error(w, "member not found", http.StatusBadRequest)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(member); err != nil {
		log.Print(err)
		return
	}
}

func (h *MemberHandler) handleDeleteMember(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
	parsedMemberId := regMemberId.FindStringSubmatch(path)[1]
	memberId, err := strconv.ParseInt(parsedMemberId, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse module id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteMember(memberId); err != nil {
		log.Print(err)
		http.Error(w, "no module found", http.StatusBadRequest)
		return
	}
}

func (h *MemberHandler) handleGetAllMembers(w http.ResponseWriter, _ *http.Request) {
	modules, err := h.db.GetAllMembers()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&modules); err != nil {
		log.Print(err)
		return
	}
}
func (h *MemberHandler) handleGetMembersByGroup(w http.ResponseWriter, r *http.Request) {
	groupIdString := r.URL.Query().Get("groupId")

	groupId, err := strconv.ParseInt(groupIdString, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "could not parse group id", http.StatusBadRequest)
		return
	}

	modules, err := h.db.GetMembersByGroup(groupId)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	header := w.Header()
	header.Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&modules); err != nil {
		log.Print(err)
		return
	}
}
