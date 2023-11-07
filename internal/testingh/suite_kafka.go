//go:build integration

package testingh

import (
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

const kafkaClientTimeout = 10 * time.Second

type KafkaSuite struct {
	ContextSuite
	createdWithinSuiteTopics []string
}

func (s *KafkaSuite) TearDownSuite() {
	if len(s.createdWithinSuiteTopics) > 0 { // Some garbage collection.
		s.deleteTopics(s.createdWithinSuiteTopics...)
	}
	s.ContextSuite.TearDownSuite()
}

func (s *KafkaSuite) KafkaBrokers() []string {
	return []string{Config.KafkaAddress}
}

func (s *KafkaSuite) RecreateTopics(topics ...string) {
	s.T().Helper()

	s.Require().NotEmpty(topics)
	s.T().Logf("recreate topics: %v", topics)

	s.deleteTopics(topics...)
	s.createTopics(topics...)
}

func (s *KafkaSuite) createTopics(topics ...string) {
	s.T().Helper()

	tc := make([]kafka.TopicConfig, 0, len(topics))
	for _, t := range topics {
		tc = append(tc, kafka.TopicConfig{
			Topic:             t,
			NumPartitions:     16,
			ReplicationFactor: 1,
		})
	}

	cResp, err := s.newKafkaClient().CreateTopics(s.SuiteCtx, &kafka.CreateTopicsRequest{Topics: tc})
	s.Require().NoError(err)
	for _, t := range topics {
		s.Require().NoError(cResp.Errors[t], "topic %q", t)
	}

	s.createdWithinSuiteTopics = append(s.createdWithinSuiteTopics, topics...)
}

func (s *KafkaSuite) deleteTopics(topics ...string) {
	s.T().Helper()

	dResp, err := s.newKafkaClient().DeleteTopics(s.SuiteCtx, &kafka.DeleteTopicsRequest{Topics: topics})
	s.Require().NoError(err)
	for _, t := range topics {
		if err, ok := dResp.Errors[t]; ok && !errors.Is(err, kafka.UnknownTopicOrPartition) {
			s.Require().NoErrorf(err, "topic %q", t)
		}
	}
}

func (s *KafkaSuite) newKafkaClient() *kafka.Client {
	return &kafka.Client{
		Addr:    kafka.TCP(s.KafkaBrokers()...),
		Timeout: kafkaClientTimeout,
	}
}
