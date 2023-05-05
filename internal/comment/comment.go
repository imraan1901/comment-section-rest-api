package comment

import (
	"context"
	"errors"
	"fmt"

	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
	"github.com/imraan1901/comment-section-rest-api/internal/processor"
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

func (s *Service) GetComment(ctx context.Context, id string) (datastructs.Comment, error) {
	fmt.Println("Retreiving a comment")
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
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
	cmt, err := s.Store.UpdateComment(ctx, id, updatedCmt)
	if err != nil {
		fmt.Println("error updating comment")
		return datastructs.Comment{}, err
	}
	return cmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.Store.DeleteComment(ctx, id)
}

func (s *Service) PostComment(ctx context.Context, cmt datastructs.Comment) (datastructs.Comment, error) {

	insertedCmt, err := s.Store.PostComment(ctx, cmt)
	if err != nil {
		return datastructs.Comment{}, err
	}

	//go func() {
	cmt, err = s.ProcessComment(ctx, insertedCmt)
	if err != nil {
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

	processor, err := processor.NewProcessor()
	if err != nil {
		return datastructs.Comment{}, err
	}

	cmt, err := processor.ProcessComment(ctx, comment)
	if err != nil {
		return datastructs.Comment{}, err
	}
	return cmt, nil
}
