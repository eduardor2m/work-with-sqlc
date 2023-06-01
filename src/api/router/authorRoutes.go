package router

import (
	"context"
	"database/sql"
	"github.com/eduardor2m/work-with-sqlc/src/infra/sqlite/bridge"
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

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS author (id INTEGER PRIMARY KEY, name TEXT NOT NULL, bio TEXT);")

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func loadAuthorRoutes(group *echo.Group) {
	authorGroup := group.Group("/author")
	authorGroup.POST("", NewAuthorHandlers().CreateAuthor)

	authorGroup.GET("/:id", NewAuthorHandlers().GetAuthor)

	authorGroup.GET("", NewAuthorHandlers().ListAuthors)

	authorGroup.DELETE("/:id", NewAuthorHandlers().DeleteAuthor)

	authorGroup.DELETE("", NewAuthorHandlers().DeleteAllAuthors)
}

type AuthorHandlers interface {
	CreateAuthor(c echo.Context) error
	GetAuthor(c echo.Context) error
	ListAuthors(c echo.Context) error
	DeleteAuthor(c echo.Context) error
	DeleteAllAuthors(c echo.Context) error
}

type authorHandlers struct{}

func NewAuthorHandlers() AuthorHandlers {
	return &authorHandlers{}
}

func (h *authorHandlers) CreateAuthor(c echo.Context) error {
	a := Author{}

	if err := c.Bind(&a); err != nil {

		return err
	}

	conn, err := dbInit()

	if err != nil {
		return err
	}

	ctx := context.Background()
	queries := bridge.New(conn)

	author, err := queries.CreateAuthor(ctx, bridge.CreateAuthorParams{
		Name: a.Name,
		Bio: sql.NullString{
			String: a.Bio,
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	return c.JSON(200, author)
}

func (h *authorHandlers) GetAuthor(c echo.Context) error {
	id := c.Param("id")

	conn, err := dbInit()

	if err != nil {
		return err
	}

	queries := bridge.New(conn)

	idToInt64, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return err
	}

	author, err := queries.GetAuthor(context.Background(), idToInt64)

	if err != nil {
		return err
	}

	return c.JSON(200, author)
}

func (h *authorHandlers) ListAuthors(c echo.Context) error {
	conn, err := dbInit()

	if err != nil {
		return err
	}

	queries := bridge.New(conn)

	authors, err := queries.ListAuthors(context.Background())

	if err != nil {
		return err
	}

	return c.JSON(200, authors)
}

func (h *authorHandlers) DeleteAuthor(c echo.Context) error {
	id := c.Param("id")

	conn, err := dbInit()

	if err != nil {
		return err
	}

	queries := bridge.New(conn)

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
}

func (h *authorHandlers) DeleteAllAuthors(c echo.Context) error {
	conn, err := dbInit()

	if err != nil {
		return err
	}

	queries := bridge.New(conn)

	ctx := context.Background()

	err = queries.DeleteAllAuthors(ctx)

	if err != nil {
		return err
	}

	return c.JSON(200, "All authors deleted")
}
