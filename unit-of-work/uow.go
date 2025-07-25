package unitofwork

import (
	"fmt"

	"github.com/akapond-katip/poc-uow-golang/repositories"
	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	Execute(fn func(rp RepoProvider) error) error
}

type RepoProvider interface {
	GetUserRepository() repositories.UserRepository
}

type RepoRegistry struct {
	tx *sqlx.Tx
	RepoProvider
}

func (rr *RepoRegistry) GetUserRepository() repositories.UserRepository {
	return repositories.NewUserRepository(rr.tx)
}

type uow struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) UnitOfWork {
	return &uow{
		db: db,
	}
}

func (u *uow) Execute(fn func(rp RepoProvider) error) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("transaction rollback failed: %v (original error: %w)", rbErr, err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("transaction commit failed: %w", err)
			}

			fmt.Println("transaction committed")
		}
	}()

	registry := &RepoRegistry{
		tx: tx,
	}
	err = fn(registry)
	return err
}
