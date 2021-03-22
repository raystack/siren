package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/odpf/siren/domain"
	"net/http"
	"net/url"
)

func UpdateAlertCredentials(service domain.AlertmanagerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var alertCredential domain.AlertCredential
		err := json.NewDecoder(r.Body).Decode(&alertCredential)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("invalid json body"))
			return
		}
		teamName := params["teamName"]
		alertCredential.TeamName = teamName
		validators := []func(*domain.AlertCredential) error{
			validateWebhooks,
			validateEntity,
			validatePagerdutyKey,
		}
		for _, v := range validators {
			if err := v(&alertCredential); err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}
		}

		err = service.Upsert(alertCredential)
		if err != nil {
			internalServerError(w, err)
			return
		}
		w.WriteHeader(201)
		return
	}
}

func validatePagerdutyKey(credential *domain.AlertCredential) error {
	if credential.PagerdutyCredentials == "" {
		return errors.New("pagerduty key cannot be empty")
	}
	return nil

}

func validateEntity(credential *domain.AlertCredential) error {
	if credential.Entity == "" {
		return errors.New("entity cannot be empty")
	}
	return nil
}

func validateWebhooks(credential *domain.AlertCredential) error {
	_, err := url.Parse(credential.SlackConfig.Critical.Webhook)
	if err != nil {
		return errors.New("slack critical webhook is not a valid url")
	}
	_, err = url.Parse(credential.SlackConfig.Warning.Webhook)
	if err != nil {
		return errors.New("slack critical webhook is not a valid url")
	}
	return nil
}
