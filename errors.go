package ep2p

import "fmt"

type topicExistsError struct {
	name string
}

func (t *topicExistsError) Error() string {
	return fmt.Sprintf("topic %s already exists", t.name)
}

func TopicAlreadyExistsError(topic string) error {
	return &topicExistsError{name: topic}
}

type topicNotFoundError struct {
	name string
}

func (t *topicNotFoundError) Error() string {
	return fmt.Sprintf("topic %s not found", t.name)
}

func TopicNotFoundError(topic string) error {
	return &topicNotFoundError{name: topic}
}

type libp2pError struct {
	err error
}

func (l *libp2pError) Error() string {
	return fmt.Sprintf("libp2p error: %v", l.err)
}

func Libp2pError(err error) error {
	return &libp2pError{err: err}
}
