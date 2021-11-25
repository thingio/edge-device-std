package operations

import (
	"fmt"
	"github.com/thingio/edge-device-std/msgbus/message"
	"github.com/thingio/edge-device-std/version"
	"strings"
)

type (
	TopicTagKey string
	TopicTags   map[TopicTagKey]string
)

const (
	TagsOffset = 2 // <OperationCategory>/<ver>/<Tags...>

	TopicTagKeyOptType    TopicTagKey = "opt_type"
	TopicTagKeyOptMode    TopicTagKey = "opt_mode"
	TopicTagKeyProtocolID TopicTagKey = "protocol_id"
	TopicTagKeyProductID  TopicTagKey = "product_id"
	TopicTagKeyDeviceID   TopicTagKey = "device_id"
	TopicTagKeyFuncID     TopicTagKey = "func_id"
	TopicTagKeyReqID      TopicTagKey = "req_id"

	TopicLevelSeparator      = "/"
	TopicMultiLevelWildcard  = "#"
	TopicSingleLevelWildcard = "+"
)

var (
	// Schemas describes all topics' forms. Every topic is formed by <ver>/<OperationCategory>/<TopicTag1>/.../<TopicTagN>.
	Schemas = map[OperationCategory][]TopicTagKey{
		OperationCategoryMeta: {TopicTagKeyOptMode, TopicTagKeyProtocolID, TopicTagKeyOptType, TopicTagKeyReqID},
		OperationCategoryData: {TopicTagKeyOptMode, TopicTagKeyProtocolID, TopicTagKeyProductID, TopicTagKeyDeviceID,
			TopicTagKeyFuncID, TopicTagKeyOptType, TopicTagKeyReqID},
	}
)

type Topic interface {
	Category() OperationCategory

	String() string

	Tags() TopicTags
	TagKeys() []TopicTagKey
	TagValues() []string
	TagValue(key TopicTagKey) (value string, ok bool)
}

type commonTopic struct {
	version  version.Version
	category OperationCategory

	tags TopicTags
}

func (c *commonTopic) Version() version.Version {
	return c.version
}

func (c *commonTopic) Category() OperationCategory {
	return c.category
}

func (c *commonTopic) String() string {
	topicVersion := c.version
	topicType := string(c.category)
	tagValues := c.TagValues()
	return strings.Join(append([]string{topicType, string(topicVersion)}, tagValues...), TopicLevelSeparator)
}

func (c *commonTopic) Tags() TopicTags {
	return c.tags
}

func (c *commonTopic) TagKeys() []TopicTagKey {
	return Schemas[c.category]
}

func (c *commonTopic) TagValues() []string {
	tagKeys := c.TagKeys()
	values := make([]string, len(tagKeys))
	for idx, topicTagKey := range tagKeys {
		values[idx] = c.tags[topicTagKey]
	}
	return values
}

func (c *commonTopic) TagValue(key TopicTagKey) (value string, ok bool) {
	tags := c.Tags()
	value, ok = tags[key]
	return
}

func ParseTopic(msg *message.Message) (Topic, error) {
	if msg == nil || msg.Topic == "" {
		return nil, fmt.Errorf("invalid message: %s", msg.String())
	}
	topic := msg.Topic

	parts := strings.Split(topic, TopicLevelSeparator)
	if len(parts) <= TagsOffset {
		return nil, fmt.Errorf("invalid topic: %s", topic)
	}
	topicCategory, topicVersion := OperationCategory(parts[0]), version.Version(parts[1])

	keys, ok := Schemas[topicCategory]
	if !ok {
		return nil, fmt.Errorf("undefined operation category: %s", topicCategory)
	}
	if len(parts)-TagsOffset != len(keys) {
		return nil, fmt.Errorf("invalid topic: %s, keys [%+v] are necessary", topic, keys)
	}

	tags := make(map[TopicTagKey]string)
	for i, key := range keys {
		tags[key] = parts[i+TagsOffset]
	}
	return &commonTopic{
		version:  topicVersion,
		category: topicCategory,
		tags:     tags,
	}, nil
}
