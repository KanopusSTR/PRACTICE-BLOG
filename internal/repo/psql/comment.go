package psql

import (
	"context"
	"server/internal/entities"
	"server/internal/repo"
	"server/pkg/myErrors"
	"time"
)

type comment struct {
	db PgxIface
}

func NewComment(db PgxIface) repo.Comment {
	return &comment{db: db}
}

func (c *comment) Add(text *string, date time.Time, authorMail string, postId int) error {
	//tx, err := c.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := c.db.Query(
		context.Background(),
		"SELECT * FROM create_comment($1, $2, $3, $4)",
		*text,
		date,
		authorMail,
		postId,
	)
	return err
}

func (c *comment) Remove(_, commentId int) error {
	//tx, err := c.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := c.db.Query(
		context.Background(),
		`SELECT * FROM delete_comment($1)`,
		commentId,
	)
	return err
}

func (c *comment) GetPostComments(postId int) ([]interface{}, error) {
	var comments []interface{}
	//tx, err := c.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	rows, err := c.db.Query(
		context.Background(),
		`SELECT * FROM read_comments($1)`,
		postId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment entities.Comment

		if err := rows.Scan(
			&comment.CommentId,
			&comment.Text,
			&comment.Date,
			&comment.AuthorMail,
			&comment.PostId,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *comment) GetPostComment(_, commentId int) (*entities.Comment, error) {
	var comm *entities.Comment
	//tx, err := c.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	rows, err := c.db.Query(
		context.Background(),
		`SELECT * FROM read_comment($1)`,
		commentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		count++
		var comment entities.Comment
		if err := rows.Scan(
			&comment.CommentId,
			&comment.Text,
			&comment.Date,
			&comment.AuthorMail,
			&comment.PostId,
		); err != nil {
			return nil, err
		}
		comm = &comment
	}
	if count != 1 {
		return nil, myErrors.CommentNotFound
	}
	return comm, nil

}
