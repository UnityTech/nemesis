package report

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/golang/glog"

	"cloud.google.com/go/pubsub"
)

// PubSubReporter is a wrapper around the pubsub.Client.
type PubSubReporter struct {
	c     *pubsub.Client
	topic string
}

// NewPubSubReporter returns a new PubSubReporter for outputting the findings of an audit
func NewPubSubReporter(project string, topic string) *PubSubReporter {
	psr := new(PubSubReporter)

	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, project)
	if err != nil {
		glog.Fatalf("Failed to create PubSub client: %v", err)
	}
	psr.c = c
	psr.topic = topic
	return psr
}

// Publish sends a list of reports to the configured PubSub topic
func (r *PubSubReporter) Publish(reports []Report) error {

	// Create a pubsub publisher workgroup
	ctx := context.Background()

	var wg sync.WaitGroup
	var errs uint64
	topic := r.c.Topic(r.topic)

	for i := 0; i < len(reports); i++ {

		// Marshal and send the report
		data, err := json.Marshal(&reports[i])
		if err != nil {
			glog.Fatalf("Failed to marshal report for pubsub: %v", err)
		}

		result := topic.Publish(ctx, &pubsub.Message{
			Data: data,
		})

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()

			_, err := res.Get(ctx)
			if err != nil {
				glog.Errorf("Failed to publish: %v", err)
				atomic.AddUint64(&errs, 1)
			}
		}(i, result)
	}

	wg.Wait()

	if errs > 0 {
		return fmt.Errorf("%d of %d reports did not publish", errs, len(reports))
	}

	return nil
}
