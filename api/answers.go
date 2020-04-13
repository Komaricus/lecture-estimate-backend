package api

import (
	"awesomeProject3/api/respond"
	"awesomeProject3/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

const (
	getAnswers = `select * from dbo.answers`

	insertAnswer = `insert into dbo.answers (
			id
			, question
			, answer
			, begin_lesson
			, date
			, discipline
			, end_lesson
			, kind_of_work
			, lecturer
			, "from"
			, will_come)
			 values (
				    :id, 
				    :question, 
				    :answer, 
					:begin_lesson,
			        :date,
			        :discipline,
				    :end_lesson,
				    :kind_of_work,
			        :lecturer,
			        :from,
			        :will_come)`
)

type Answer struct {
	ID          string `json:"id" db:"id"`
	Question    string `json:"question" db:"question"`
	Answer      int    `json:"answer" db:"answer"`
	BeginLesson string `json:"beginLesson" db:"begin_lesson"`
	Date        string `json:"date" db:"date"`
	Discipline  string `json:"discipline" db:"discipline"`
	EndLesson   string `json:"endLesson" db:"end_lesson"`
	KindOfWork  string `json:"kindOfWork" db:"kind_of_work"`
	Lecturer    string `json:"lecturer" db:"lecturer"`
	From        string `json:"from" db:"from"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	WillCome    bool   `json:"willCome" db:"will_come"`
}

// Answers list
//
// swagger:response answersList
type GetAnswersResponseBody []*Answer

// Answers swagger:route GET /answers Answers GetAnswers
//
// Gets answers from answers table
//
// Responses:
//   200: answersList
//   500: errorResponse
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	answers, err := getAnswersFromDB(ctx)
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, answers)
}

// Error while processing request
//
// swagger:parameters PostAnswer
type postAnswerRequest struct {
	// in: body
	Body RequestPostAnswer
}

type RequestPostAnswer struct {
	Question    string `json:"question"`
	Answer      int    `json:"answer"`
	BeginLesson string `json:"beginLesson"`
	Date        string `json:"date"`
	Discipline  string `json:"discipline"`
	EndLesson   string `json:"endLesson"`
	KindOfWork  string `json:"kindOfWork"`
	Lecturer    string `json:"lecturer"`
	From        string `json:"from"`
	WillCome    bool   `json:"willCome"`
}

func convertIncomingAnswerToDBAnswer(in *RequestPostAnswer) *Answer {
	return &Answer{
		ID:          uuid.New().String(),
		Question:    in.Question,
		Answer:      in.Answer,
		BeginLesson: in.BeginLesson,
		Date:        in.Date,
		Discipline:  in.Discipline,
		EndLesson:   in.EndLesson,
		KindOfWork:  in.KindOfWork,
		Lecturer:    in.Lecturer,
		From:        in.From,
		WillCome:    in.WillCome,
	}
}

// Answers swagger:route POST /answer Answers PostAnswer
//
// Adds answers to answers table
//
// Responses:
//   200: statusOkResponse
//   400: badRequestResponse
//   500: errorResponse
func PostAnswer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	var rpa RequestPostAnswer
	if err := json.NewDecoder(r.Body).Decode(&rpa); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err := insertAnswerToDB(ctx, convertIncomingAnswerToDBAnswer(&rpa)); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, nil)
}

// Error while processing request
//
// swagger:parameters PostAnswers
type postAnswersRequest struct {
	// in: body
	Body []RequestPostAnswer
}

// Answers swagger:route POST /answers Answers PostAnswers
//
// Adds many answers to answers table
//
// Responses:
//   200: statusOkResponse
//   400: badRequestResponse
//   500: errorResponse
func PostAnswers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	var rpa []RequestPostAnswer
	if err := json.NewDecoder(r.Body).Decode(&rpa); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	for _, ans := range rpa {
		if err := insertAnswerToDB(ctx, convertIncomingAnswerToDBAnswer(&ans)); err != nil {
			respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}
	}

	respond.WithJSON(ctx, w, http.StatusOK, nil)
}

func getAnswersFromDB(ctx context.Context) ([]Answer, error) {
	db := models.GetDB()

	answers := make([]Answer, 0)
	if err := db.SelectContext(ctx, &answers, getAnswers); err != nil {
		return nil, err
	}

	return answers, nil
}

func insertAnswerToDB(ctx context.Context, ans *Answer) error {
	db := models.GetDB()

	if _, err := db.NamedExecContext(ctx, insertAnswer, &ans); err != nil {
		return err
	}

	return nil
}
