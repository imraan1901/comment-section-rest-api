package db

// This file in the db package queries the database and 
// returns the result to the business layer comment/comment.go

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/imraan1901/comment-section-rest-api/internal/comment"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Slug.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComment(
	ctx context.Context,
	uuid string,
) (comment.Comment, error) {

	var cmtRow CommentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, slug, body, author
		 FROM comments
		 WHERE id=$1`,
		uuid,
	)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching comment by uuid: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}
