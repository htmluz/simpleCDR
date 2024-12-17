package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"radiusgo/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func EventStream(q *services.CallQueue) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				count := q.GetQueueSize()
				calls := q.GetAllCalls()
				jcalls, e := json.Marshal(calls)
				if e != nil {
					fmt.Printf("error serializing json calls %v", e)
				}
				fmt.Fprintf(w, "data: {\"count\": %d, \"calls\": %s}\n\n", count, jcalls)

				err := w.Flush()
				if err != nil {
					break
				}
				time.Sleep(2000 * time.Millisecond)
			}
		}))
		return nil
	}

}
