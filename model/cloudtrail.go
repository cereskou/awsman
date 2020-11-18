package model

import "gorm.io/gorm"

//CloudTrailEvent -
type CloudTrailEvent struct {
	gorm.Model
	EventID         string `json:"id" gorm:"index:cloudtrail_id,unique"`
	AccessKeyID     string `json:"accessKey"`
	CloudTrailEvent string `json:"event"`
	EventName       string `json:"name"`
	EventSource     string `json:"source"`
	EventTime       string `json:"eventtime"`
	Username        string `json:"username"`
}

//CloudTrailEventResource -
type CloudTrailEventResource struct {
	gorm.Model
	EventID string `json:"id" gorm:"index:cloudtrail_resource,unique"`
	Name    string `json:"resourceName" gorm:"index:cloudtrail_resource,unique"`
	Type    string `json:"resourceType"`
}
