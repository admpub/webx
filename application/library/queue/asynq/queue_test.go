package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/admpub/pp"
	"github.com/hibiken/asynq"
)

func TestQueue(t *testing.T) {
	a, err := New(nil)
	if err != nil {
		panic(err)
	}
	a.HandleFunc(`test:1`, func(ctx context.Context, t *asynq.Task) error {
		data := map[string]interface{}{}
		err := json.Unmarshal(t.Payload(), &data)
		if err != nil {
			return err
		}
		id := data["user_id"]
		fmt.Printf("[a] Send Welcome Email to User %d\n", id)
		return nil
	})
	defer a.Close()
	go a.StartWorker()

	a1, err := New(nil)
	if err != nil {
		panic(err)
	}
	a1.HandleFunc(`test:1`, func(ctx context.Context, t *asynq.Task) error {
		data := map[string]interface{}{}
		err := json.Unmarshal(t.Payload(), &data)
		if err != nil {
			return err
		}
		id := data["user_id"]
		fmt.Printf("[a1] Send Welcome Email to User %d\n", id)
		return nil
	})
	defer a1.Close()
	go a1.StartWorker()

	a2, err := New(nil)
	if err != nil {
		panic(err)
	}
	a.HandleFunc(`test:2`, func(ctx context.Context, t *asynq.Task) error {
		data := map[string]interface{}{}
		err := json.Unmarshal(t.Payload(), &data)
		if err != nil {
			return err
		}
		id := data["user_id"]
		fmt.Printf("[a2] Send Welcome Email to User %d\n", id)
		return nil
	})
	defer a2.Close()
	go a2.StartWorker()

	c, err := New(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 20; i++ {
		var typename string
		if i%2 == 0 {
			typename = `test:1`
		} else {
			typename = `test:2`
		}
		data := map[string]interface{}{`user_id`: i}
		payload, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		result, err := c.Send(asynq.NewTask(typename, payload))
		if err != nil {
			fmt.Println(err)
		}
		pp.Println(result)
		time.Sleep(1 * time.Second)
	}
	c.Close()
}
