package bmcluster

import "fmt"

type Port struct {
	BmId  int
	Index int
}

type Message struct {
	From Port
	To   Port
}

func (p Port) String() string {
	return fmt.Sprintf("Port{BmId: %d, Index: %d}", p.BmId, p.Index)
}

func (m Message) String() string {
	return fmt.Sprintf("Message from %v to %v", m.From, m.To)
}

func (c *Cluster) GetMessages() []Message {
	messages := make([]Message, 0)
	for _, peerI := range c.Peers {
		for i, input := range peerI.Inputs {
		osearch:
			for _, peerO := range c.Peers {
				for o, output := range peerO.Outputs {
					if input == output {
						messages = append(messages, Message{
							To:   Port{BmId: int(peerI.PeerId), Index: i},
							From: Port{BmId: int(peerO.PeerId), Index: o},
						})
						break osearch
					}
				}
			}
		}
	}
	return messages
}
