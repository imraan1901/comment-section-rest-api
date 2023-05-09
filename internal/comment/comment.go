package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
	"github.com/imraan1901/comment-section-rest-api/internal/processor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

var (
	ErrFetchingComment = errors.New("failed to fetch comment by id")
	ErrNotImplemented  = errors.New("not implemented")
)

// Store - this interface defines all of the methods
// that our service needs in order to operate
// It is of this type if it implements these functions
type Store interface {
	GetComment(context.Context, string) (datastructs.Comment, error)
	PostComment(context.Context, datastructs.Comment) (datastructs.Comment, error)
	UpdateComment(context.Context, string, datastructs.Comment) (datastructs.Comment, error)
	DeleteComment(context.Context, string) error
}

// Service - is the struct in which
// all of our logic wil be built on
// All services of this type are also
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new
// service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

// name is the Tracer name used to identify this instrumentation library.
const name = "comment"

func (s *Service) GetComment(ctx context.Context, id string) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "GetComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	fmt.Println("Retreiving a comment")
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		fmt.Println(err)
		return datastructs.Comment{}, ErrFetchingComment
	}

	return cmt, nil
}

func (s *Service) UpdateComment(
	ctx context.Context,
	id string,
	updatedCmt datastructs.Comment,
) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "UpdateComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	cmt, err := s.Store.UpdateComment(ctx, id, updatedCmt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		fmt.Println("error updating comment")
		return datastructs.Comment{}, err
	}
	return cmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "DeleteComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	return s.Store.DeleteComment(ctx, id)
}

func (s *Service) PostComment(ctx context.Context, cmt datastructs.Comment) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "PostComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	insertedCmt, err := s.Store.PostComment(ctx, cmt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, err
	}

	//go func() {
	cmt, err = s.ProcessComment(ctx, insertedCmt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		fmt.Println("error processing comment: %w", err)
		return datastructs.Comment{}, err
	}
	s.UpdateComment(ctx, cmt.ID, cmt)
	//}()

	return insertedCmt, nil
}

// In a production environment never let user input be executed with eval
// for demo purposes only
func (s *Service) ProcessComment(ctx context.Context, comment datastructs.Comment) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "ProcessComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	processor, err := processor.NewProcessor()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, err
	}

	cmt, err := processor.ProcessComment(ctx, comment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, err
	}
	return cmt, nil
}
