package psql

import (
	"context"
	"server/internal/entities"
	"server/internal/repo"
	"server/pkg/myErrors"
)

type user struct {
	db PgxIface
}

func NewUser(db PgxIface) repo.User {
	return &user{db: db}
}

func (u *user) Add(user entities.User) error {
	//tx, err := u.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := u.db.Query(
		context.Background(),
		"SELECT * FROM create_user($1, $2, $3)",
		user.Name,
		user.Mail,
		user.Password,
	)
	return err
}

func (u *user) Get(mail string) (*entities.User, error) {
	var user *entities.User
	//tx, err := u.db.Begin(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	rows, err := u.db.Query(
		context.Background(),
		`SELECT * FROM read_user($1)`,
		mail,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0

	for rows.Next() {
		count++
		var user2 entities.User
		if err := rows.Scan(
			&user2.Name,
			&user2.Mail,
			&user2.Password,
		); err != nil {
			return nil, err
		}
		user = &user2
	}
	if count != 1 {
		return nil, myErrors.UserNotFound
	}
	return user, nil
}
