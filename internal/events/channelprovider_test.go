package events

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/vseinstrumentiru/lego/multilog"
)

func TestPubSub(t *testing.T) {
	ass := assert.New(t)

	c, err := ProvideChannel(channelArgs{
		Logger: multilog.New(0),
	})
	ass.NoError(err)

	ch, err := c.Subscribe(context.Background(), "test.me")
	ass.NoError(err)

	var counter int
	wg := &sync.WaitGroup{}
	go func() {
		for msg := range ch {
			if strings.HasPrefix(string(msg.Payload), "test message") {
				counter++
			}
			msg.Ack()
			wg.Done()
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		err := c.Publish("test.me", message.NewMessage(uuid.New().String(), []byte(fmt.Sprintf("test message %v", i))))
		ass.NoError(err)
	}

	wg.Wait()
	ass.NoError(c.Close())
	ass.Equal(10, counter)
}
