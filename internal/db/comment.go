package db

// This file in the db package queries the database and
// returns the result to the business layer comment/comment.go

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imraan1901/comment-section-rest-api/internal/datastructs"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(c CommentRow) datastructs.Comment {
	return datastructs.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Body:   c.Slug.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetComment(
	ctx context.Context,
	uuid string,
) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "GetComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, fmt.Errorf("error fetching comment by uuid: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}

func (d *Database) PostComment(ctx context.Context, cmt datastructs.Comment) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "PostComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	cmt.ID = uuid.NewV4().String()

	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}
	rows, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO comments
		(id, slug, author, body)
		VALUES
		(:id, :slug, :author, :body)`,
		postRow,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}
	if err := rows.Close(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return cmt, nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "DeleteComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments where id=$1`,
		id,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		fmt.Errorf("failed to delete comment from database: %w", err)
	}
	return nil
}

func (d *Database) UpdateComment(
	ctx context.Context,
	id string,
	cmt datastructs.Comment,
) (datastructs.Comment, error) {

	startTime := time.Now()
	_, span := otel.Tracer(name).Start(ctx, "UpdateComment", tr.WithTimestamp(startTime))
	defer span.End(tr.WithTimestamp(time.Now()))

	cmtRow := CommentRow{
		ID:     id,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE comments SET
		slug = :slug,
		author = :author,
		body = :body
		WHERE id = :id`,
		cmtRow,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}

	if err := rows.Close(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return datastructs.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	}

	return convertCommentRowToComment(cmtRow), nil

}
