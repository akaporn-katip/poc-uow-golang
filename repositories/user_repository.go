package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db sqlx.ExtContext
}

type UserEntity struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}

func NewUserRepository(db sqlx.ExtContext) UserRepository {
	return UserRepository{db: db}
}

func (ur UserRepository) NewUser(name string, email string) error {
	ctx := context.Background()
	rs, err := ur.db.ExecContext(ctx, "insert into users values ($1, $2)", name, email)

	if err != nil {
		err = fmt.Errorf("error %w", err)
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		err = fmt.Errorf("error %w", err)
		return err
	}

	fmt.Printf("Insert ID == %v\n", id)
	return nil
}

func (ur UserRepository) Lists() ([]UserEntity, error) {
	ctx := context.Background()
	users := []UserEntity{}
	rows, err := ur.db.QueryxContext(ctx, "select * from users")

	if err != nil {
		err = fmt.Errorf("error %w", err)
		return nil, err
	}

	for rows.Next() {
		user := UserEntity{}
		err := rows.StructScan(&user)
		if err != nil {
			err = fmt.Errorf("error %w", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// func (ur UserRepository) ListUser()
