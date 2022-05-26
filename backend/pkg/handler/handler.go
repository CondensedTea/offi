package handler

import (
	"errors"
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/core"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const appName = "offi-backend"

type Handler struct {
	app          *fiber.App
	sessionStore *session.Store

	core *core.Core
}

func New(c *core.Core) *Handler {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		GETOnly:      true,
	})

	prometheus := fiberprometheus.New(appName)
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)
	app.Use(logger.New())

	handler := &Handler{
		app:          app,
		sessionStore: session.New(),
		core:         c,
	}

	handler.app.Get("/match/:matchId", handler.GetMatch)
	handler.app.Get("/log/:logId", handler.GetLog)
	handler.app.Get("/recruitment/:type/:id", handler.GetRecPost)
	handler.app.Get("", handler.Debug)

	return handler
}

func (h Handler) Run() error {
	return h.app.Listen(":8080")
}

func (h Handler) GetMatch(ctx *fiber.Ctx) error {
	matchId, err := ctx.ParamsInt("matchId")
	if err != nil {
		return err
	}
	if matchId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "match id is required")
	}
	logs, err := h.core.GetLogs(matchId)
	if errors.Is(err, cache.ErrCached) {
		return fiber.NewError(fiber.StatusTooEarly, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get logs: %v", err))
	}

	sess, err := h.sessionStore.Get(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get session: %v", err))
	}

	if err = sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save session: %v", err))
	}

	views, err := h.core.CountViews("match", matchId, sess.Fresh())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get view count for match: %v", err))
	}
	return ctx.JSON(fiber.Map{"logs": logs, "views": views})
}

func (h Handler) GetLog(ctx *fiber.Ctx) error {
	logId, err := ctx.ParamsInt("logId")
	if err != nil {
		return err
	}
	if logId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "log id is required")
	}
	match, err := h.core.GetMatch(logId)
	if err == redis.Nil {
		return fiber.NewError(fiber.StatusNotFound, "this log does not have linked match yet")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get match: %v", err))
	}

	sess, err := h.sessionStore.Get(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get session: %v", err))
	}
	if err = sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save session: %v", err))
	}
	views, err := h.core.CountViews("logs", logId, sess.Fresh())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update view count for log: %v", err))
	}
	return ctx.JSON(fiber.Map{"match": match, "views": views})
}

func (h Handler) GetRecPost(ctx *fiber.Ctx) error {
	postType := ctx.Params("type")
	id := ctx.Params("id")

	if postType == "" || id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you must specify post type and entity ID"))
	}

	var (
		entry *cache.Entry
		err   error
	)
	switch postType {
	case "player":
		entry, err = h.core.GetPlayerRecruitmentStatus(id)
	case "team":
		entry, err = h.core.GetTeamRecruitmentStatus(id)
	default:
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("unknown post type: %s", postType))
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get player recruitment status: %v", err))
	}
	return ctx.JSON(fiber.Map{"status": entry})
}

func (h Handler) Debug(ctx *fiber.Ctx) error {
	hashKey := ctx.Params("hashKey")

	keys, err := h.core.GetKeys(hashKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get keys: %v", err))
	}
	return ctx.JSON(fiber.Map{"keys": keys, "count": len(keys)})
}
