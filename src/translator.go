package main

import (
	"encoding/json"
	translator "github.com/gilang-as/google-translate"
	"log"
)

type Language struct {
	DidYouMean bool   `json:"did_you_mean"`
	Iso        string `json:"iso"`
}

type Text struct {
	AutoCorrected bool   `json:"auto_corrected"`
	Value         string `json:"value"`
	DidYouMean    bool   `json:"did_you_mean"`
}

type From struct {
	Language `json:"language"`
	Text     `json:"text"`
}

type TranslateResponse struct {
	Text          string `json:"text"`
	Pronunciation string `json:"pronunciation"`
	From          `json:"from"`
}

func translateTo(text string, to string) (TranslateResponse, error) {

	value := translator.Translate{
		Text: text,
		To:   to,
	}
	translated, err := translator.Translator(value)
	if err != nil {
		log.Println(err)
		return TranslateResponse{}, err
	} else {
		prettyJSON, err := json.MarshalIndent(translated, "", "\t")
		var response TranslateResponse
		err = json.Unmarshal(prettyJSON, &response)
		if err != nil {
			log.Println(err)
			return TranslateResponse{}, err
		}
		log.Println("translate result: " + string(prettyJSON))
		return response, nil
	}
}
