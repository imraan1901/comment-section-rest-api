package processor

import (
	"context"
	"fmt"

	//python3 "github.com/go-python/cpy3"
	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
)

type ProcessStatus int

const (
	Failed      ProcessStatus = iota - 2 // -1
	Processing                = iota     // 0
	Processed                 = iota     // 1
	UnProcessed               = iota     // 2
)

// Returns any type interface for our processor 
// so any module can call this code
type PService struct {
	processor interface{}
}

// Any code can call the process comment function when a NewProcessor is made
func NewProcessor() (*PService, error) {

	var any interface{}
	return &PService{
		processor: any,
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

	//defer python3.Py_Finalize()
	//python3.Py_Initialize()
	

	var pcmt ProcessedComment
	var err error
	if err != nil {
		fmt.Println("Failed to connect")
		return datastructs.Comment{}, err
	}

	//var processed interface{}

	//processed = python3.PyRun_SimpleString("print('hello world')")
	pcmt.Processed_Author = fmt.Sprintf("%v", "processed")
	pcmt.Processed_Body = fmt.Sprintf("%v", "processed")
	pcmt.Processed_Slug = fmt.Sprintf("%v", "processed")

	newCmt := datastructs.Comment{
		ID:     cmt.ID,
		Slug:   pcmt.Processed_Slug,
		Body:   pcmt.Processed_Body,
		Author: pcmt.Processed_Author,
	}


	/*
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
	*/

	return newCmt, nil
}
