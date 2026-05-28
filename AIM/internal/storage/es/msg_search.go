package es

import (
	"bytes"
	"encoding/json"
	"fmt"

	"AIM/internal/storage/mysql/model"
)

const Index = "im_message"

func AddMsgToES(msg *model.Message) error {
	data, _ := json.Marshal(msg)
	_, err := Client.Index(Index, bytes.NewReader(data))
	return err
}

func SearchMsg(keyword string) ([]model.Message, error) {
	query := fmt.Sprintf(`{"query":{"match":{"content":"%s"}}}`, keyword)
	res, err := Client.Search(
		Client.Search.WithIndex(Index),
		Client.Search.WithBody(bytes.NewBufferString(query)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r struct {
		Hits struct {
			Hits []struct {
				Source model.Message `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	json.NewDecoder(res.Body).Decode(&r)

	msgs := make([]model.Message, 0, len(r.Hits.Hits))
	for _, v := range r.Hits.Hits {
		msgs = append(msgs, v.Source)
	}
	return msgs, nil
}
