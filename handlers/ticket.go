package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"github.com/vector-ops/ticketeer/models"
)

type TicketHandler struct {
	repository models.TicketRepository
}

func NewTicketHandler(router fiber.Router, repo models.TicketRepository) *TicketHandler {
	handler := &TicketHandler{
		repository: repo,
	}

	router.Get("/", handler.GetMany)
	router.Get("/:ticketId", handler.GetOne)
	router.Post("/", handler.CreateOne)
	router.Put("/:ticketId", handler.UpdateOne)
	router.Delete("/:ticketId", handler.DeleteOne)
	router.Post("/validate", handler.ValidateOne)

	return handler
}

func (h *TicketHandler) GetMany(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	userId := uint(ctx.Locals("userId").(float64))

	tickets, err := h.repository.GetMany(context, userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data":    tickets,
	})
}

func (h *TicketHandler) GetOne(ctx *fiber.Ctx) error {
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	userId := uint(ctx.Locals("userId").(float64))

	ticket, err := h.repository.GetOne(context, uint(ticketId), userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	var QRCode []byte

	QRCode, err = qrcode.Encode(fmt.Sprintf("ticketId:%v, ownerId:%v", ticketId, userId), qrcode.Medium, 256)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data": &fiber.Map{
			"ticket": ticket,
			"qrcode": QRCode,
		},
	})
}

func (h *TicketHandler) CreateOne(ctx *fiber.Ctx) error {
	ticket := &models.Ticket{}
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	if err := ctx.BodyParser(ticket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	userId := uint(ctx.Locals("userId").(float64))

	ticket, err := h.repository.CreateOne(context, ticket, userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "Ticket created",
		"data":    ticket,
	})
}

func (h *TicketHandler) ValidateOne(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	validateTicket := &models.ValidateTicket{}

	if err := ctx.BodyParser(validateTicket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	validateData := make(map[string]interface{})
	validateData["entered"] = true

	ticket, err := h.repository.UpdateOne(context, validateTicket.TicketId, validateData, validateTicket.UserId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "Ticket valid",
		"data":    ticket,
	})
}

func (h *TicketHandler) UpdateOne(ctx *fiber.Ctx) error {
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	userId := uint(ctx.Locals("userId").(float64))

	updateTicket := make(map[string]interface{})
	if err := ctx.BodyParser(&updateTicket); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	ticket, err := h.repository.UpdateOne(context, uint(ticketId), updateTicket, userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "Updated ticket",
		"data":    ticket,
	})
}
func (h *TicketHandler) DeleteOne(ctx *fiber.Ctx) error {
	ticketId, _ := strconv.Atoi(ctx.Params("ticketId"))

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	userId := uint(ctx.Locals("userId").(float64))

	err := h.repository.DeleteOne(context, uint(ticketId), userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
