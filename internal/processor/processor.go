package processor

import (
	"context"
	"fmt"

	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
	python3 "github.com/go-python/cpy3"
)

type ProcessStatus int

const (
	Failed      ProcessStatus = iota - 2 // -1
	Processing                = iota     // 0
	Processed                 = iota     // 1
	UnProcessed               = iota     // 2
)

type PService struct {
	client python3.PyObject
}

func NewProcessor() (*PService, error) {

	client = python.Py_Initialize()

	if !python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}

	//processor, err := client.
	if err != nil {
		return nil, err
	}

	return &PService{
		client: processor,
	}, nil
}

type ProcessedComment struct {
	ID               string
	Processed_Slug   string
	Processed_Body   string
	Processed_Author string
}

// Takes in a PService and updates its values with the valid uuid
// In a production environment never let user input be executed with eval
// for demo purposes only
func (w *PService) ProcessComment(
	ctx context.Context,
	cmt datastructs.Comment) (datastructs.Comment, error) {

	var pcmt ProcessedComment
	var err error
	if err != nil {
		fmt.Println("Failed to connect")
		return datastructs.Comment{}, err
	}

	var processed interface{}

	processed, err = w.client.Eval("pi")
	if err != nil {
		fmt.Println("Failed to process data")
		return datastructs.Comment{}, err
	}
	pcmt.Processed_Author = fmt.Sprintf("%v", processed)

	processed, err = w.client.Eval("pi")
	if err != nil {
		fmt.Println("Failed to process data")
		return datastructs.Comment{}, err
	}
	pcmt.Processed_Body = fmt.Sprintf("%v", processed)

	processed, err = w.client.Eval("pi")
	if err != nil {
		fmt.Println("Failed to process data")
		return datastructs.Comment{}, err
	}
	pcmt.Processed_Slug = fmt.Sprintf("%v", processed)

	newCmt := datastructs.Comment{
		ID:     cmt.ID,
		Slug:   pcmt.Processed_Slug,
		Body:   pcmt.Processed_Body,
		Author: pcmt.Processed_Author,
	}

	return newCmt, nil
}
