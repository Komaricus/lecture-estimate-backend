package api

import (
	"awesomeProject3/api/respond"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

type Lesson struct {
	Auditorium  string `json:"auditorium"`
	BeginLesson string `json:"beginLesson"`
	Building    string `json:"building"`
	Date        string `json:"date"`
	DayOfWeek   int    `json:"dayOfWeek"`
	Discipline  string `json:"discipline"`
	EndLesson   string `json:"endLesson"`
	KindOfWork  string `json:"kindOfWork"`
	Lecturer    string `json:"lecturer"`
}

// Error while processing request
//
// swagger:parameters GetLessons
type pathGetLessons struct {
	// group / student
	//
	// required: true
	request_type string

	// group / student id
	//
	// required: true
	id string

	// start date example - 2020.01.20
	//
	// required: true
	start_date string

	// end date example - 2020.01.26
	//
	// required: true
	end_date string
}

// Lessons
//
// swagger:response lessonList
type GetLessonsResponseBody []*Lesson

// Lessons swagger:route GET /lessons Lessons GetLessons
//
// Get lessons by filter criteria
//
// Responses:
//   200: lessonList
//   400: badRequestResponse
//   500: errorResponse
func GetLessons(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	requestType := r.URL.Query().Get("request_type")
	if requestType == "" {
		err := errors.New("invalid params:request_type")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("request_type", requestType)

	id := r.URL.Query().Get("id")
	if id == "" {
		err := errors.New("invalid params:id")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("id", id)

	startDate := r.URL.Query().Get("start_date")
	if startDate == "" {
		err := errors.New("invalid params:start_date")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("start_date", startDate)

	endDate := r.URL.Query().Get("end_date")
	if endDate == "" {
		err := errors.New("invalid params:end_date")
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}
	logger.WithField("end_date", endDate)

	lessons, err := makeLessonsRequest(requestType, id, startDate, endDate)
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, lessons)
}

func makeLessonsRequest(requestType, id, startDate, endDate string) ([]Lesson, error) {
	url := fmt.Sprintf("https://ruz.hse.ru/api/schedule/%v/%v?start=%v&finish=%v&lng=1", requestType, id, startDate, endDate)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	lessons := make([]Lesson, 0)
	if err := json.NewDecoder(resp.Body).Decode(&lessons); err != nil {
		return nil, err
	}

	return lessons, nil
}
