package proto

type RequestOrResponse struct {
	Size int32
	M
}
type Request struct {
	APIKey         int16
	APIVersion     int16
	CorrelationID  int32
	ClientID       string
	RequestMessage RequestMessage
}
type Response struct {
	CorrelationID   int32
	ResponseMessage ResponseMessage
}
type ResponseMessage M
type MessageSet []OffsetMessage
type OffsetMessage struct {
	Offset       int64
	SizedMessage SizedMessage
}
type SizedMessage struct {
	Size       int32
	CRCMessage CRCMessage
}
type CRCMessage struct {
	CRC     uint32
	Message Message
}
type Message struct {
	MagicByte  int8
	Attributes int8
	Key        []byte
	Value      []byte
}
type TopicMetadataRequest []string
type TopicMetadataResponse struct {
	Brokers        []Broker
	TopicMetadatas []TopicMetadata
}
type Broker struct {
	NodeID int32
	Host   string
	Port   int32
}
type TopicMetadata struct {
	ErrorCode          ErrorCode
	TopicName          string
	PartitionMetadatas []PartitionMetadata
}
type PartitionMetadata struct {
	ErrorCode   ErrorCode
	PartitionID int32
	Leader      int32
	Replicas    []int32
	ISR         []int32
}
type ProduceRequest struct {
	RequiredAcks       int16
	Timeout            int32
	MessageSetInTopics []MessageSetInTopic
}
type MessageSetInTopic struct {
	TopicName              string
	MessageSetInPartitions []MessageSetInPartition
}
type MessageSetInPartition struct {
	Partition  int32
	MessageSet MessageSet
}
type ProduceResponse []OffsetInTopic
type OffsetInTopic struct {
	TopicName          string
	OffsetInPartitions []OffsetInPartition
}
type OffsetInPartition struct {
	Partition int32
	ErrorCode ErrorCode
	Offset    int64
}
type FetchRequest struct {
	ReplicaID           int32
	MaxWaitTime         int32
	MinBytes            int32
	FetchOffsetInTopics []FetchOffsetInTopic
}
type FetchOffsetInTopic struct {
	TopicName               string
	FetchOffsetInPartitions []FetchOffsetInPartition
}
type FetchOffsetInPartition struct {
	Partition   int32
	FetchOffset int64
	MaxBytes    int32
}
type FetchResponse []FetchMessageSetInTopic
type FetchMessageSetInTopic struct {
	TopicName                   string
	FetchMessageSetInPartitions []FetchMessageSetInPartition
}
type FetchMessageSetInPartition struct {
	Partition           int32
	ErrorCode           ErrorCode
	HighwaterMarkOffset int64
	MessageSet          MessageSet
}
type OffsetRequest struct {
	ReplicaID    int32
	TimeInTopics []TimeInTopic
}
type TimeInTopic struct {
	TopicName        string
	TimeInPartitions []TimeInPartition
}
type TimeInPartition struct {
	Partition          int32
	Time               int64
	MaxNumberOfOffsets int32
}
type OffsetResponse []OffsetsInTopic
type OffsetsInTopic struct {
	TopicName           string
	OffsetsInPartitions []OffsetsInPartition
}
type OffsetsInPartition struct {
	Partition int32
	ErrorCode ErrorCode
	Offsets   []int64
}
type ConsumerMetadataRequest string
type ConsumerMetadataResponse struct {
	ErrorCode       ErrorCode
	CoordinatorID   int32
	CoordinatorHost string
	CoordinatorPort int32
}
type OffsetCommitRequestV0 struct {
	ConsumerGroupID        string
	OffsetCommitInTopicV0s []OffsetCommitInTopicV0
}
type OffsetCommitInTopicV0 struct {
	TopicName                  string
	OffsetCommitInPartitionV0s []OffsetCommitInPartitionV0
}
type OffsetCommitInPartitionV0 struct {
	Partition int32
	Offset    int64
	Metadata  string
}
type OffsetCommitRequestV1 struct {
	ConsumerGroupID           string
	ConsumerGroupGenerationID int32
	ConsumerID                string
	OffsetCommitInTopicV1s    []OffsetCommitInTopicV1
}
type OffsetCommitInTopicV1 struct {
	TopicName                  string
	OffsetCommitInPartitionV1s []OffsetCommitInPartitionV1
}
type OffsetCommitInPartitionV1 struct {
	Partition int32
	Offset    int64
	TimeStamp int64
	Metadata  string
}
type OffsetCommitRequestV2 struct {
	ConsumerGroup             string
	ConsumerGroupGenerationID int32
	ConsumerID                string
	RetentionTime             int64
	OffsetCommitInTopicV2s    []OffsetCommitInTopicV2
}
type OffsetCommitInTopicV2 struct {
	TopicName                  string
	OffsetCommitInPartitionV2s []OffsetCommitInPartitionV2
}
type OffsetCommitInPartitionV2 struct {
	Partition int32
	Offset    int64
	Metadata  string
}
type OffsetCommitResponse []ErrorInTopic
type ErrorInTopic struct {
	TopicName         string
	ErrorInPartitions []ErrorInPartition
}
type ErrorInPartition struct {
	Partition int32
	ErrorCode ErrorCode
}
type OffsetFetchRequestV0 struct {
	ConsumerGroup     string
	PartitionInTopics []PartitionInTopic
}
type PartitionInTopic struct {
	TopicName  string
	Partitions []int32
}
type OffsetFetchRequestV1 struct {
	ConsumerGroup     string
	PartitionInTopics []PartitionInTopic
}
type OffsetFetchRequestV2 struct {
	ConsumerGroup     string
	PartitionInTopics []PartitionInTopic
}
type OffsetFetchResponse []OffsetMetadataInTopic
type OffsetMetadataInTopic struct {
	TopicName                  string
	OffsetMetadataInPartitions []OffsetMetadataInPartition
}
type OffsetMetadataInPartition struct {
	Partition int32
	Offset    int64
	Metadata  string
	ErrorCode ErrorCode
}
type ErrorCode int16
