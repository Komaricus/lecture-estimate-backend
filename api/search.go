package api

import (
	"awesomeProject3/api/respond"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

const timeOut = 60 * time.Second

type Group struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// Error while processing request
//
// swagger:parameters GetGroups
type pathSearchGroup struct {
	// search string example - БИВ
	//
	// required: true
	search string
}

// Groups
//
// swagger:response groupList
type GetGroupsResponseBody []*Group

// Groups swagger:route GET /groups Groups GetGroups
//
// Get groups by search filter
//
// Responses:
//   200: groupList
//   400: badRequestResponse
//   500: errorResponse
func GetGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	search := r.URL.Query().Get("search")
	if search == "" {
		err := errors.New("invalid params:search")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("search", search)

	groups, err := MakeSearchRequest(search, "group")
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, groups)
}

// Error while processing request
//
// swagger:parameters GetStudents
type pathSearchStudent struct {
	// search string example - Александр ...
	//
	// required: true
	search string
}

// Students
//
// swagger:response studentList
type GetStudentsResponseBody []*Group

// Students swagger:route GET /students Students GetStudents
//
// Get students by search filter
//
// Responses:
//   200: studentList
//   400: badRequestResponse
//   500: errorResponse
func GetStudents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	search := r.URL.Query().Get("search")
	if search == "" {
		err := errors.New("invalid params:search")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("search", search)

	students, err := MakeSearchRequest(search, "student")
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, students)
}

func MakeSearchRequest(search, searchType string) ([]Group, error) {
	search = strings.ReplaceAll(search, " ", "%20")
	resp, err := http.Get(fmt.Sprintf("https://ruz.hse.ru/api/search?term=%s&type=%v", search, searchType))
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0)
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, err
	}

	return groups, nil
}
