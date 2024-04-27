package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Invoices []struct {
// 	Sent bool
// 	Id   int
// }

type Queue struct {
	items []string
}

type FormData struct {
	Items []string `json:"invoices"`
}

func NewQueue() *Queue {
	return &Queue{
		items: []string{},
	}
}

func (q *Queue) Enqueue(item string) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *Queue) DequeueSpecific(value interface{}) (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	for i, item := range q.items {
		if item == value {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return item, nil
		}
	}

	return nil, errors.New("invoice not found in queue")
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Length() int {
	return len(q.items)
}

func main() {
	q := NewQueue()
	r := gin.Default()

	r.GET("/invoices", func(c *gin.Context) {
		if q.IsEmpty() {
			c.JSON(404, gin.H{
				"error": "No items in queue.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"invoices": q.items,
		})
	})

	r.DELETE("/invoices", func(c *gin.Context) {
		var formData FormData
		if err := c.BindJSON(&formData); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid data",
			})
			return
		}

		if len(formData.Items) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})

			return
		}

		if q.IsEmpty() {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "No invoices in queue",
			})
			return
		}

		for _, value := range formData.Items {
			if item, err := q.DequeueSpecific(value); err == nil {
				fmt.Printf("Dequeued: %v\n", item)
			} else {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message":            "Processed",
			"processed-invoices": formData.Items,
		})

	})

	r.POST("/invoices", func(c *gin.Context) {
		var formData FormData
		if err := c.BindJSON(&formData); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid data",
			})
			return
		}

		if len(formData.Items) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})

			return
		}

		stringSet := make(map[string]struct{})

		for _, str := range q.items {
			stringSet[str] = struct{}{}
		}

		for _, newStr := range formData.Items {
			if _, exists := stringSet[newStr]; !exists {
				q.Enqueue(newStr)
				stringSet[newStr] = struct{}{}
			}
		}

		c.JSON(200, gin.H{
			"message":           "Received invoices",
			"received_invoices": formData.Items,
		})
	})

	r.Run("0.0.0.0:9090")
}
