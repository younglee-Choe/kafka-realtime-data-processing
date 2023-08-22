package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"main/src/goDotEnvVariable"
	"main/src/structures"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers                  = []string{goDotEnvVariable.GoDotEnvVariable("KAFKA_1"), goDotEnvVariable.GoDotEnvVariable("KAFKA_2"), goDotEnvVariable.GoDotEnvVariable("KAFKA_3")}
	topic        goka.Stream = "leele-topic"
	topicRekeyed goka.Stream = "leele-topic-rekeyed"
	group        goka.Group  = "leele-group"

	tmc                         *goka.TopicManagerConfig
	producerSize                int = 4
	defaultPartitionChannelSize int = 4
)

// This codec allows marshalling (encode) and unmarshalling (decode) the block struct(or produce struct) to and from the group table
type stateDateCodec struct{}
type blockDataCodec struct{}

func init() {
	tmc = goka.NewTopicManagerConfig()
	tmc.Table.Replication = 3
	tmc.Stream.Replication = 3
}

// Encodes structures.StateData into []byte
func (jc *stateDateCodec) Encode(value interface{}) ([]byte, error) {
	if _, isState := value.(*structures.StateData); !isState {
		return nil, fmt.Errorf("Codec requires value *structures.StateData, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a structures.StateData from []byte to it's go representation
func (jc *stateDateCodec) Decode(data []byte) (interface{}, error) {
	var (
		c   structures.StateData
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling structures.StateData: %v", err)
	}
	return &c, nil
}

// Encodes structures.BlockData into []byte
func (jc *blockDataCodec) Encode(value interface{}) ([]byte, error) {
	if _, isBlock := value.(*structures.BlockData); !isBlock {
		return nil, fmt.Errorf("Codec requires value *structures.BlockData, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes a structures.BlockData from []byte to it's go representation
func (jc *blockDataCodec) Decode(data []byte) (interface{}, error) {
	var (
		c   structures.BlockData
		err error
	)
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling structures.BlockData: %v", err)
	}
	return &c, nil
}

// process() re-keys; Producer's ID -> offset(index)
func process(ctx goka.Context, msg interface{}) {
	key := ctx.Offset()
	ctx.Loopback(strconv.Itoa(int(key)), msg)
}

// loopProcess() uses the context to update the state of the StateData structures's data
func loopProcess(ctx goka.Context, msg interface{}) {
	var s *structures.StateData
	if val := ctx.Value(); val != nil {
		s = val.(*structures.StateData)
	} else {
		s = new(structures.StateData)
	}

	if s.Counter < producerSize {
		// The counter is incremented when data with the same key comes in
		s.Counter++

		// Append encoded data(source data) sent by producer
		s.EncodedData = append(s.EncodedData, []byte(msg.(*structures.BlockData).SourceData))

		if s.Counter == 1 { // Initialize Length field of StateData struct
			s.Length = msg.(*structures.BlockData).Length
			s.IsLength = false
		} else if !s.IsLength && s.Counter > 1 {
			if s.Length != msg.(*structures.BlockData).Length {
				// If the previous state and the current value are not the same, save the current value and continue
				s.Length = msg.(*structures.BlockData).Length
			} else {
				// Otherwise (if equal), save the current value and do not change it
				s.IsLength = true
			}
		}
	} else {
		s.Counter = 0
	}
	// SetValue() updates the value of the key in the group table.
	// It stores the value in the local cache and sends the update to the Kafka topic representing the group table
	ctx.SetValue(s)

	// When all the data of the same offset comes together, try data processing
	if s.Counter == producerSize {
		// ❕Conditional real-time state-based data processing
	}
}

// Write a processor that consumes data from Kafka
func runProcessor(initialized chan struct{}) {
	g := goka.DefineGroup(group,
		goka.Input(topic, new(blockDataCodec), process), // function for receiving messages(stream) from Kafka
		goka.Loop(new(blockDataCodec), loopProcess),     // re-key using status
		goka.Output(topicRekeyed, new(codec.String)),    // push messages(stream) to Kafka
		goka.Persist(new(stateDateCodec)),               // required for stateful-based data(table) processing
	)
	p, err := goka.NewProcessor(brokers,
		g,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		panic(err)
	}

	close(initialized)

	if err = p.Run(context.Background()); err != nil {
		log.Printf("Error running processor: %v", err)
	}
}

// Writing a view to query the user table
func runView(initialized chan struct{}) {
	<-initialized

	view, err := goka.NewView(brokers,
		goka.GroupTable(group),
		new(blockDataCodec),
	)
	if err != nil {
		panic(err)
	}

	view.Run(context.Background())
}

func main() {
	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		log.Fatalf("Error creating topic manager: %v", err)
	}
	defer tm.Close()
	err = tm.EnsureStreamExists(string(topic), defaultPartitionChannelSize)
	if err != nil {
		log.Printf("Error creating kafka topic %s: %v", topic, err)
	}

	// When this example is run the first time, wait for creation of all internal topics (this is done
	// by goka.NewProcessor)
	initialized := make(chan struct{})

	go runProcessor(initialized)
	runView(initialized)
}
