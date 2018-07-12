package message

import (
	"errors"
	"log"
	"mqtt/utils"
	"strings"
)

// Topic :
type Topic struct {
	Name string
	Desc string
}

// Message :
type Message struct {
	ID        int64
	TopicName string
	Data      []byte
}

// Sub : subscription
type Sub struct {
	ClientID    string
	TopicFilter string
}

var (
	topicMap = make(map[string]*Topic)
	subs     []*Sub
	messages []*Message
)

// NewTopic : create a new topic
func NewTopic(topic *Topic) error {
	log.Println("creating topic => ", topic.Name)
	_, ok := topicMap[topic.Name]
	if ok {
		return errors.New("duplicate topic")
	}
	topicMap[topic.Name] = topic
	return nil
}

// DestroyTopic : destroy a existing topic
func DestroyTopic(name string) {
	log.Println("destroying topic => ", name)
	delete(topicMap, name)
}

// MatchTopic :
func MatchTopic(pattern string) {
	log.Println("matching topic pattern => ", pattern)
}

// NewSub : create a new subscription
func NewSub(s *Sub) error {
	if strings.EqualFold(s.ClientID, utils.Blank) {
		return errors.New("client_id is empty")
	}
	if strings.EqualFold(s.TopicFilter, utils.Blank) {
		return errors.New("no topic filter found")
	}
	subs = append(subs, s)
	return nil
}

// DeleteSubs : delete subscriptions
func DeleteSubs(clientID string, topicFilters []string) {
	for i := len(subs) - 1; i >= 0; i-- {
		for _, v := range topicFilters {
			if strings.EqualFold(clientID, subs[i].ClientID) &&
				strings.EqualFold(v, subs[i].TopicFilter) {
				// remove the element subs[i]
				subs = append(subs[:i], subs[i+1:]...)
			}
		}
	}
}

// GetSubs : get subscriptions
func GetSubs() []*Sub {
	return subs
}

// PutMessage : put a new message into the messages
func PutMessage(message *Message) error {
	_, ok := topicMap[message.TopicName]
	if !ok {
		return errors.New("topic doesn't exist")
	}
	var lastIndex int64
	if len(messages) > 0 {
		lastIndex = messages[len(messages)-1].ID
	}
	lastIndex++
	message.ID = lastIndex
	messages = append(messages, message)
	return nil
}

// PopMessage : pop out a message from the messages
func PopMessage() *Message {
	if len(messages) == 0 {
		return nil
	}
	var message []*Message
	copy(message, messages[len(messages)-1:])
	messages = append(messages[len(messages)-1:])
	return message[0]
}
