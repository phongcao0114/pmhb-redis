package models

import "time"

const (
	// EventPHTest const
	EventPHTest string = "test"
)

type (
	// EventMessage structure
	EventMessage struct {
		EventName    string      `json:"event_name" bson:"event_name"`
		QueueID      string      `json:"queue_id" bson:"queue_id"`
		Domain       string      `json:"domain" bson:"domain"`
		EventVersion int         `json:"event_version" bson:"event_version"`
		Segment      Segment     `json:"segment" bson:"segment"`
		Meta         Meta        `json:"meta" bson:"meta"`
		Payload      interface{} `json:"payload,omitempty" bson:"payload"`
		Topic        string      `bson:"topic"`
	}
	// Meta structure
	Meta struct {
		RequestID        string    `json:"request_id" bson:"request_id"`
		RequestAppID     string    `json:"request_app_id" bson:"request_app_id"`
		RequestTimeStamp time.Time `json:"request_timestamp" bson:"request_time_stamp"`
	}
	// Segment structure
	Segment struct {
		CurrentPartition int `json:"current_partition" bson:"current_partition"`
		TotalPartition   int `json:"total_partition" bson:"total_partition"`
	}
)
