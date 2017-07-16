package esa

import (
	"bytes"
	"encoding/json"
)

const (
	CategoryURL = "/v1/teams"
)

type CategoryService struct {
	client *Client
}

type CategoryBatchMoveReq struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type CategoryBatchMoveRes struct {
	Count int    `json:"count"`
	From  string `json:"from"`
	To    string `json:"to"`
}

func (s *CategoryService) BatchMove(teamName string, categoryFrom string, categoryTo string) error {
	batchMoveURL := CategoryURL + "/" + teamName + "/categories/batch_move"
	batchMoveReq := CategoryBatchMoveReq{ From: categoryFrom, To: categoryTo }

	var batchMoveRes CategoryBatchMoveRes
	var data []byte
	var err error
	if data, err = json.Marshal(batchMoveReq); err != nil {
		return err
	}

	res, err := s.client.post(batchMoveURL, "application/json", bytes.NewReader(data), &batchMoveRes)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
