package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"radiusgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func EventStream(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			calls, count := services.GetActiveCalls()
			callsjson, e := json.Marshal(calls)
			if e != nil {
				log.Printf("erro serializando o json %v", e)
				return
			}
			msg := fmt.Sprintf("data: {\"calls\": %s, \"count\": %v}\n\n", callsjson, count)
			fmt.Fprintf(w, "data: {%s}\n\n", msg)

			err := w.Flush()
			if err != nil {
				fmt.Printf("closing http conn, %v", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}))

	return nil
}
