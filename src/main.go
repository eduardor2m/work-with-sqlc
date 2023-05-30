package main

import (
	"context"
	"database/sql"
	db "github.com/eduardor2m/work-with-sqlc/db/sqlc"
	"github.com/labstack/echo/v4"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Author struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func dbInit() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "./authors.db")

	if err != nil {
		return nil, err
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS authors (id INTEGER PRIMARY KEY, name TEXT NOT NULL, bio TEXT);")

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func dbCreateAuthor(conn *sql.DB, name string, bio string) (*db.Author, error) {
	ctx := context.Background()
	queries := db.New(conn)

	author, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: name,
		Bio: sql.NullString{
			String: bio,
			Valid:  true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &author, nil
}

func server() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, world!")
	})

	api := e.Group("/api")

	authors := api.Group("/authors")

	authors.GET("", func(c echo.Context) error {
		return c.String(200, "Hello, authors!")
	})

	authors.POST("", func(c echo.Context) error {

		a := Author{}

		if err := c.Bind(&a); err != nil {

			return err
		}

		conn, err := dbInit()

		if err != nil {
			return err
		}

		author, err := dbCreateAuthor(conn, a.Name, a.Bio)

		if err != nil {
			return err
		}

		return c.JSON(200, author)

	})

	authors.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")

		conn, err := dbInit()

		if err != nil {
			return err
		}

		queries := db.New(conn)

		idToInt64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			return err
		}

		author, err := queries.GetAuthor(context.Background(), idToInt64)

		if err != nil {
			return err
		}

		return c.JSON(200, author)
	})

	authors.GET("", func(c echo.Context) error {
		conn, err := dbInit()

		if err != nil {
			return err
		}

		queries := db.New(conn)

		authors, err := queries.ListAuthors(context.Background())

		if err != nil {
			return err
		}

		return c.JSON(200, authors)

	})

	authors.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")

		conn, err := dbInit()

		if err != nil {
			return err
		}

		queries := db.New(conn)

		idToInt64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			return err
		}

		ctx := context.Background()

		err = queries.DeleteAuthor(ctx, idToInt64)

		if err != nil {
			return err
		}

		return c.JSON(200, "Author deleted")
	})

	authors.DELETE("", func(c echo.Context) error {
		conn, err := dbInit()

		if err != nil {
			return err
		}

		queries := db.New(conn)

		ctx := context.Background()

		err = queries.DeleteAllAuthors(ctx)

		if err != nil {
			return err
		}

		return c.JSON(200, "All authors deleted")
	})

	e.Logger.Fatal(e.Start(":8080"))

}

func main() {
	_, err := dbInit()

	if err != nil {
		panic(err)
	}

	server()
}
