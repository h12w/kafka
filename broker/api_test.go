package broker

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"h12.me/realtest/kafka"
	//"h12.me/wipro"
)

func TestTopicMetadata(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	b := New(DefaultConfig().WithAddr(leaderAddr))
	defer b.Close()
	key, value := "test key", "test value"
	produceMessage(t, b, topic, partition, key, value)
	messages := fetchMessage(t, b, topic, partition, 0)
	if len(messages) != 1 {
		t.Fatalf("expect 1 message but got %v", messages)
	}
	if m := messages[0]; m[0] != key || m[1] != value {
		t.Fatalf("expect [%s %s] but got %v", key, value, m)
	}
}

func TestProduceSnappy(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	partitionCount := 2
	partition := int32(1)
	topic := "topic1"
	err = k.NewTopic(topic, partitionCount)
	if err != nil {
		t.Fatal(err)
	}
	defer k.DeleteTopic(topic)
	leaderAddr := getLeader(t, k, topic, partition)

	conn, err := net.Dial("tcp", leaderAddr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write([]byte{0x00, 0x00, 0x00, 0xc3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x00, 0x01, 0x00, 0x00, 0x03, 0xe8, 0x00, 0x00, 0x00, 0x01, 0x00, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x93, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3f, 0xb9, 0xa5, 0xb2, 0xea, 0x00, 0x02, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x31, 0x50, 0x00, 0x00, 0x19, 0x01, 0x64, 0x19, 0x5a, 0x0b, 0x91, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x0b, 0x2d, 0x2d, 0x2d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x2d, 0x2d, 0x2d, 0x19, 0x24, 0x14, 0x00, 0x1f, 0x18, 0x8c, 0x9b, 0x8d, 0x15, 0x25, 0x04, 0x11, 0x78, 0x3e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c, 0x74, 0xef, 0x27, 0xd6, 0x00, 0x02, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x2e, 0x42, 0x00, 0x00, 0x19, 0x01, 0x3c, 0x13, 0x50, 0x6e, 0x03, 0xa6, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x05, 0x63, 0x01, 0x01, 0x19, 0x1e, 0x14, 0x00, 0x17, 0x30, 0xb5, 0x42, 0x70, 0x15, 0x1f, 0x24, 0x09, 0x2d, 0x2d, 0x2d, 0x65, 0x6e, 0x64, 0x2d, 0x2d, 0x2d})
	if err != nil {
		t.Fatal(err)
	}

	//b := New(DefaultConfig().WithAddr(leaderAddr))
	//defer b.Close()
	/*
		var w wipro.Writer
		ms := MessageSet{
			{SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
				Attributes: 0,
				Key:        nil,
				Value:      []byte("hello"),
			}}}},
		}
		ms.Marshal(&w)
		compressedValue := encodeSnappy(w.B)
		fmt.Println(w.B)
		resp, err := b.Produce(topic, partition, MessageSet{
			{
				SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
					Attributes: 2,
					Key:        nil,
					Value:      compressedValue,
				}}}},
		})
	*/
	resp := &ProduceResponse{}
	r := &Response{
		ResponseMessage: resp,
	}
	if err := r.Receive(conn); err != nil {
		t.Fatal(err)
	}
	if len(*resp) != 1 || len((*resp)[0].OffsetInPartitions) != 1 {
		t.Fatal("expect 1 resp")
	}

	for _, topic := range *resp {
		for _, p := range topic.OffsetInPartitions {
			if p.ErrorCode.HasError() {
				t.Fatal(p.ErrorCode)
			}
		}
	}
	t.Log(*resp)
}

func TestGroupCoordinator(t *testing.T) {
	t.Parallel()
	k, err := kafka.New()
	if err != nil {
		t.Fatal(err)
	}
	group := kafka.RandomGroup()
	coord := getCoord(t, k, group)
	fmt.Println(group, coord)
}

func getTopicMetadata(t *testing.T, k *kafka.Cluster, topic string) *TopicMetadataResponse {
	b := New(DefaultConfig().WithAddr(k.AnyBroker()))
	defer b.Close()
	respMsg, err := b.TopicMetadata(topic)
	if err != nil {
		t.Fatal(err)
	}
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

func getCoord(t *testing.T, k *kafka.Cluster, group string) string {
	reqMsg := GroupCoordinatorRequest(group)
	req := &Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &GroupCoordinatorResponse{}
	resp := &Response{ResponseMessage: respMsg}
	conn, err := net.Dial("tcp", k.AnyBroker())
	if err != nil {
		t.Fatal(err)
	}
	sendReceive(t, conn, req, resp)
	if respMsg.ErrorCode.HasError() {
		t.Fatal(respMsg.ErrorCode)
	}
	return respMsg.Broker.Addr()
}

func produceMessage(t *testing.T, b *B, topic string, partition int32, key, value string) {
	respMsg, err := b.Produce(topic, partition, []OffsetMessage{
		{
			SizedMessage: SizedMessage{CRCMessage: CRCMessage{Message: Message{
				Key:   []byte(key),
				Value: []byte(value),
			}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
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

func fetchMessage(t *testing.T, b *B, topic string, partition int32, offset int64) [][2]string {
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
	if err := b.Do(req, &respMsg); err != nil {
		t.Fatal(err)
	}
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

func sendReceive(t *testing.T, conn net.Conn, req *Request, resp *Response) {
	if err := req.Send(conn); err != nil {
		t.Fatal(t)
	}
	if err := resp.Receive(conn); err != nil {
		t.Fatal(err)
	}
	if resp.CorrelationID != req.CorrelationID {
		t.Fatalf("correlation id: expect %d but got %d", req.CorrelationID, resp.CorrelationID)
	}
}
