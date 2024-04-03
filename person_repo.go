package main

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

type personRepoContainer struct {
	db *sqlx.DB
}

type PersonRepo interface {
	Insert(p Person) error
}

func (p *personRepoContainer) Insert(person Person) error {
	_, e := p.db.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", person)
	return e
}

func BuildPersonRepo(db *dbContainer) (PersonRepo, error) {
	return &personRepoContainer{
		db: db.db,
	}, nil
}

func BuildInputPersonHandler(repo PersonRepo) HandlerDef {
	return HandlerDef{
		Path: "/person",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e := repo.Insert(Person{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@doe.com",
			})

			if e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
		}),
	}
}

var person = fx.Module("person",
	fx.Provide(
		BuildPersonRepo,
		// fx.Annotate(BuildInputPersonHandler, fx.ResultTags(`group:"handlers"`)),
	),
)
