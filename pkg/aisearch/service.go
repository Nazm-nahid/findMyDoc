package aisearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type DeepSeekService interface {
	SuggestSpecialist(symptoms string) (string, error)
}

type deepSeek struct {
	apiKey string
}

func NewDeepSeekService(apiKey string) DeepSeekService {
	return &deepSeek{apiKey}
}

func (d *deepSeek) SuggestSpecialist(symptoms string) (string, error) {
	// Construct the request payload
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful medical assistant. Based on symptoms, suggest the most suitable type of medical specialist."},
			{"role": "user", "content": "Symptoms: " + symptoms + ". Suggest a doctor category like Cardiologist, Neurologist, etc. Only output the category."},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("DeepSeek API error response:", string(bodyBytes))
		return "", errors.New("DeepSeek API error: " + resp.Status)
	}

	// Parse the response
	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(data, &res); err != nil {
		log.Println("Failed to parse DeepSeek response:", string(data))
		return "", err
	}

	// Return the suggestion
	if len(res.Choices) > 0 {
		return res.Choices[0].Message.Content, nil
	}

	return "", errors.New("No suggestion returned from DeepSeek")
}
