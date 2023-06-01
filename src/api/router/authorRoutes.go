package router

import (
	"context"
	"database/sql"
	db "github.com/eduardor2m/work-with-sqlc/db/sqlc"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
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

func loadAuthorRoutes(group *echo.Group) {
	authorGroup := group.Group("/author")
	authorGroup.POST("", func(c echo.Context) error {

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

	authorGroup.GET("/:id", func(c echo.Context) error {
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

	authorGroup.GET("", func(c echo.Context) error {
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

	authorGroup.DELETE("/:id", func(c echo.Context) error {
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

	authorGroup.DELETE("", func(c echo.Context) error {
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
}
