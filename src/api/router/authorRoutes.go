package router

import (
	"context"
	"database/sql"
	"github.com/eduardor2m/work-with-sqlc/src/infra/sqlite"
	"github.com/eduardor2m/work-with-sqlc/src/infra/sqlite/bridge"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"strconv"
)

type Author struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
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

	conn, err := sqlite.GetConnection()

	if err != nil {
		return err
	}

	defer sqlite.CloseConnection(conn)

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

	conn, err := sqlite.GetConnection()

	if err != nil {
		return err
	}

	defer sqlite.CloseConnection(conn)

	queries := bridge.New(conn)

	idToInt64, err := strconv.ParseInt(id, 10, 64)

	idToInt32 := int32(idToInt64)

	if err != nil {
		return err
	}

	author, err := queries.GetAuthor(context.Background(), idToInt32)

	if err != nil {
		return err
	}

	return c.JSON(200, author)
}

func (h *authorHandlers) ListAuthors(c echo.Context) error {
	conn, err := sqlite.GetConnection()

	if err != nil {
		return err
	}

	defer sqlite.CloseConnection(conn)

	queries := bridge.New(conn)

	authors, err := queries.ListAuthors(context.Background())

	if err != nil {
		return err
	}

	return c.JSON(200, authors)
}

func (h *authorHandlers) DeleteAuthor(c echo.Context) error {
	id := c.Param("id")

	conn, err := sqlite.GetConnection()

	if err != nil {
		return err
	}

	defer sqlite.CloseConnection(conn)

	queries := bridge.New(conn)

	idToInt64, err := strconv.ParseInt(id, 10, 64)

	idToInt32 := int32(idToInt64)

	if err != nil {
		return err
	}

	ctx := context.Background()

	err = queries.DeleteAuthor(ctx, idToInt32)

	if err != nil {
		return err
	}

	return c.JSON(200, "Author deleted")
}

func (h *authorHandlers) DeleteAllAuthors(c echo.Context) error {
	conn, err := sqlite.GetConnection()

	if err != nil {
		return err
	}

	defer sqlite.CloseConnection(conn)

	queries := bridge.New(conn)

	ctx := context.Background()

	err = queries.DeleteAllAuthors(ctx)

	if err != nil {
		return err
	}

	return c.JSON(200, "All authors deleted")
}
