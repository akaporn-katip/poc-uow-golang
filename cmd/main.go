package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	unitofwork "github.com/akapond-katip/poc-uow-golang/unit-of-work"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("Connect Failed")
		os.Exit(1)
	}

	uow := unitofwork.New(db)
	e := echo.New()
	db.Exec(`CREATE TABLE IF NOT EXISTS users(
		name varchar,
		email varchar
	)`)

	e.POST("/", func(c echo.Context) error {
		return uow.Execute(func(rp unitofwork.RepoProvider) error {
			var req = User{}
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			userRepos := rp.GetUserRepository()

			fmt.Printf("Insert %v %v\n", req.Name, req.Email)

			err := userRepos.NewUser(req.Name, req.Email)

			if err != nil {
				fmt.Printf("%w", err.Error())
				return c.String(400, "Bad request")
			}

			if req.Name == "akaporn" {

				c.String(400, "Rollback")
				return fmt.Errorf("Rollback")
			}

			return c.String(http.StatusOK, "OK")
		})
	})

	e.GET("/", func(c echo.Context) error {
		return uow.Execute(func(rp unitofwork.RepoProvider) error {
			userRepos := rp.GetUserRepository()

			lists, err := userRepos.Lists()

			if err != nil {
				err = fmt.Errorf("%w", err)
				return c.String(400, "Bad request")
			}

			users := []User{}

			for _, u := range lists {
				users = append(users, User{
					Name:  u.Name,
					Email: u.Email,
				})
			}

			return c.JSON(http.StatusOK, users)
		})
	})

	e.Start(":3333")

	// uow.GetUserRepository()
}
