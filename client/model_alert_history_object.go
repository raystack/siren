/*
 * Siren.
 *
 * Documentation of our Siren API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client
import (
	"time"
)

type AlertHistoryObject struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	Id int32 `json:"id,omitempty"`
	Level string `json:"level,omitempty"`
	MetricName string `json:"metric_name,omitempty"`
	MetricValue string `json:"metric_value,omitempty"`
	Name string `json:"name,omitempty"`
	TemplateId string `json:"template_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}