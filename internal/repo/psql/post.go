package psql

import (
	"context"
	"server/internal/entities"
	"server/internal/repo"
	"server/pkg/myErrors"
	"time"
)

type post struct {
	db PgxIface
}

func NewPost(db PgxIface) repo.Post {
	return &post{db: db}
}

func (p *post) Add(header, body *string, date time.Time, authorMail string) error {
	//tx, err := p.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := p.db.Query(
		context.Background(),
		"SELECT * FROM create_post($1, $2, $3, $4)",
		*header,
		*body,
		date,
		authorMail,
	)
	return err
}

func (p *post) Remove(postId int) error {
	//tx, err := p.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := p.db.Query(
		context.Background(),
		"SELECT * FROM delete_post($1)",
		postId,
	)
	return err
}

func (p *post) Update(postId int, header, body *string) error {
	//tx, err := p.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := p.db.Query(
		context.Background(),
		"SELECT * FROM update_post($1, $2, $3)",
		postId,
		*header,
		*body,
	)
	return err
}

func (p *post) GetPost(id int) (*entities.Post, error) {
	var post *entities.Post
	//tx, err := p.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	rows, err := p.db.Query(
		context.Background(),
		`SELECT * FROM read_post($1)`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		count++
		var post2 entities.Post
		if err := rows.Scan(
			&post2.Id,
			&post2.Header,
			&post2.Body,
			&post2.Date,
			&post2.AuthorMail,
		); err != nil {
			return nil, err
		}
		post = &post2
	}
	if count != 1 {
		return nil, myErrors.PostNotFound
	}
	return post, nil
}

func (p *post) GetPosts() ([]interface{}, error) {
	var posts []interface{}
	//tx, err := p.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	rows, err := p.db.Query(
		context.Background(),
		`SELECT * FROM read_posts()`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post entities.Post
		if err := rows.Scan(
			&post.Id,
			&post.Header,
			&post.Body,
			&post.Date,
			&post.AuthorMail,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
