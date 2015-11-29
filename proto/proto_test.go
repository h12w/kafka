package proto

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"h12.me/realtest/kafka"
)

func TestTopicMetadata(t *testing.T) {
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	topic, err := k.NewRandomTopic(partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	respMsg := getTopicMetadata(t, k, topic)
	meta := &respMsg.TopicMetadatas[0]
	if len(meta.PartitionMetadatas) != partitionCount {
		t.Fatalf("partition count: expect %d but got %d", partitionCount, len(meta.PartitionMetadatas))
	}
	for _, pMeta := range meta.PartitionMetadatas {
		if pMeta.ErrorCode != NoError {
			t.Fatal(pMeta.ErrorCode)
		}
	}
}

func TestProduceFetch(t *testing.T) {
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	partition := int32(1)
	topic, err := k.NewRandomTopic(partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	leaderAddr := getLeader(t, k, topic, partition)
	conn, err := net.Dial("tcp", leaderAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	key, value := "test key", "test value"
	produceMessage(t, conn, topic, partition, key, value)
	messages := fetchMessage(t, conn, topic, partition, 0)
	if len(messages) != 1 {
		t.Fatalf("expect 1 message but got %v", messages)
	}
	if m := messages[0]; m[0] != key || m[1] != value {
		t.Fatalf("expect [%s %s] but got %v", key, value, m)
	}
}

func sendReceive(t *testing.T, conn net.Conn, req *Request, resp *Response) {
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(resp)
	}
	if resp.CorrelationID != req.CorrelationID {
		t.Fatalf("correlation id: expect %d but got %d", req.CorrelationID, resp.CorrelationID)
	}
}

func getTopicMetadata(t *testing.T, k *kafka.Cluster, topic string) *TopicMetadataResponse {
	conn, err := net.Dial("tcp", k.AnyBroker())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	req := &Request{
		CorrelationID:  rand.Int31(),
		RequestMessage: &TopicMetadataRequest{topic},
	}
	respMsg := &TopicMetadataResponse{}
	resp := &Response{ResponseMessage: respMsg}
	sendReceive(t, conn, req, resp)
	brokers := k.Brokers()
	for i := range brokers {
		if respMsg.Brokers[i].Addr() != brokers[i] {
			t.Fatalf("broker: expect %s but got %s", brokers[i], respMsg.Brokers[i].Addr())
		}
	}
	if len(respMsg.TopicMetadatas) != 1 {
		t.Fatalf("len(TopicMetadatas): expect 1 but got %d", len(respMsg.TopicMetadatas))
	}
	meta := &respMsg.TopicMetadatas[0]
	if meta.ErrorCode != NoError {
		t.Fatal(meta.ErrorCode)
	}
	if meta.TopicName != topic {
		t.Fatalf("topic: expect %s but got %s", topic, meta.TopicName)
	}
	return respMsg
}

func getLeader(t *testing.T, k *kafka.Cluster, topic string, partitionID int32) string {
	metaResp := getTopicMetadata(t, k, topic)
	meta := &metaResp.TopicMetadatas[0]
	leaderAddr := ""
	for _, partition := range meta.PartitionMetadatas {
		if partition.PartitionID == partitionID {
			for _, broker := range metaResp.Brokers {
				if broker.NodeID == partition.Leader {
					leaderAddr = broker.Addr()
				}
			}
		}
	}
	if leaderAddr == "" {
		t.Fatalf("fail to find leader in topic %s partition %d", topic, partitionID)
	}
	return leaderAddr
}

func produceMessage(t *testing.T, conn net.Conn, topic string, partition int32, key, value string) {
	req := &Request{
		CorrelationID: rand.Int31(),
		RequestMessage: &ProduceRequest{
			RequiredAcks: AckLocal,
			Timeout:      int32(time.Second / time.Millisecond),
			MessageSetInTopics: []MessageSetInTopic{
				{
					TopicName: topic,
					MessageSetInPartitions: []MessageSetInPartition{
						{
							Partition: partition,
							MessageSet: []OffsetMessage{
								{
									SizedMessage: SizedMessage{CRCMessage: CRCMessage{
										Message: Message{
											Key:   []byte(key),
											Value: []byte(value),
										},
									}}},
							},
						},
					},
				},
			},
		},
	}
	respMsg := &ProduceResponse{}
	resp := &Response{ResponseMessage: respMsg}
	sendReceive(t, conn, req, resp)
	ok := false
	for _, topicResp := range *respMsg {
		if topicResp.TopicName == topic {
			for _, partitionResp := range topicResp.OffsetInPartitions {
				if partitionResp.Partition == partition {
					if ok = (partitionResp.ErrorCode == NoError); !ok {
						t.Fatal(partitionResp.ErrorCode)
					}
				}
			}
		}
	}
	if !ok {
		t.Fatal("produce failed")
	}
}

func fetchMessage(t *testing.T, conn net.Conn, topic string, partition int32, offset int64) [][2]string {
	req := &Request{
		CorrelationID: rand.Int31(),
		RequestMessage: &FetchRequest{
			ReplicaID:   -1,
			MaxWaitTime: int32(time.Second / time.Millisecond),
			MinBytes:    1,
			FetchOffsetInTopics: []FetchOffsetInTopic{
				{
					TopicName: topic,
					FetchOffsetInPartitions: []FetchOffsetInPartition{
						{
							Partition:   partition,
							FetchOffset: offset,
							MaxBytes:    1024 * 1024,
						},
					},
				},
			},
		},
	}
	respMsg := FetchResponse{}
	resp := &Response{ResponseMessage: &respMsg}
	sendReceive(t, conn, req, resp)
	var result [][2]string
	for _, t := range respMsg {
		if t.TopicName == topic {
			for _, p := range t.FetchMessageSetInPartitions {
				if p.Partition == partition {
					if p.ErrorCode == NoError {
						for _, msg := range p.MessageSet {
							m := &msg.SizedMessage.CRCMessage.Message
							result = append(result, [2]string{string(m.Key), string(m.Value)})
						}
					}
				}
			}
		}
	}
	return result
}
