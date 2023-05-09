package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

type Response struct {
	Message string
}

type CommentService interface {
	PostComment(context.Context, datastructs.Comment) (datastructs.Comment, error)
	GetComment(ctx context.Context, ID string) (datastructs.Comment, error)
	UpdateComment(ctx context.Context, ID string, newCmt datastructs.Comment) (datastructs.Comment, error)
	DeleteComment(ctx context.Context, ID string) error
}

// Validate input from http request
type PostCommentRequest struct {
	Slug   string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func convertPostCommentRequestToComment(c PostCommentRequest) datastructs.Comment {
	return datastructs.Comment{
		Slug:   c.Slug,
		Author: c.Author,
		Body:   c.Body,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	ctx, span := otel.Tracer(name).Start(r.Context(), "PostComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	var cmt PostCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	validate := validator.New()
	err := validate.Struct(cmt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "not a valid comment", http.StatusBadRequest)
		return
	}

	convertedComment := convertPostCommentRequestToComment(cmt)
	postedComment, err := h.Service.PostComment(ctx, convertedComment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Print(err)
		return
	}

	if err := json.NewEncoder(w).Encode(postedComment); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		panic(err)
	}

}
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	ctx, span := otel.Tracer(name).Start(r.Context(), "GetComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.GetComment(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		panic(err)
	}

}
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	ctx, span := otel.Tracer(name).Start(r.Context(), "UpdateComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cmt datastructs.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	cmt, err := h.Service.UpdateComment(ctx, id, cmt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		panic(err)
	}

}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	ctx, span := otel.Tracer(name).Start(r.Context(), "DeleteComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteComment(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted"}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		panic(err)
	}

}
