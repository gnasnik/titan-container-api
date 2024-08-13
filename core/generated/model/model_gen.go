// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package model

import (
	"time"
)

// deployments
type Deployment struct {
	ID         string    `db:"id" json:"id"`
	Owner      string    `db:"owner" json:"owner"`
	Name       string    `db:"name" json:"name"`
	State      int32     `db:"state" json:"state"`
	Type       int32     `db:"type" json:"type"`
	Authority  bool      `db:"authority" json:"authority"`
	Version    string    `db:"version" json:"version"`
	Balance    float64   `db:"balance" json:"balance"`
	Cost       float64   `db:"cost" json:"cost"`
	ProviderID string    `db:"provider_id" json:"provider_id"`
	Expiration time.Time `db:"expiration" json:"expiration"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type LocationCn struct {
	ID        int64     `db:"id" json:"id"`
	Ip        string    `db:"ip" json:"ip"`
	Continent string    `db:"continent" json:"continent"`
	Country   string    `db:"country" json:"country"`
	Province  string    `db:"province" json:"province"`
	City      string    `db:"city" json:"city"`
	Longitude string    `db:"longitude" json:"longitude"`
	AreaCode  string    `db:"area_code" json:"area_code"`
	Latitude  string    `db:"latitude" json:"latitude"`
	Isp       string    `db:"isp" json:"isp"`
	ZipCode   string    `db:"zip_code" json:"zip_code"`
	Elevation string    `db:"elevation" json:"elevation"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type LocationEn struct {
	ID        int64     `db:"id" json:"id"`
	Ip        string    `db:"ip" json:"ip"`
	Continent string    `db:"continent" json:"continent"`
	Country   string    `db:"country" json:"country"`
	Province  string    `db:"province" json:"province"`
	City      string    `db:"city" json:"city"`
	Longitude string    `db:"longitude" json:"longitude"`
	AreaCode  string    `db:"area_code" json:"area_code"`
	Latitude  string    `db:"latitude" json:"latitude"`
	Isp       string    `db:"isp" json:"isp"`
	ZipCode   string    `db:"zip_code" json:"zip_code"`
	Elevation string    `db:"elevation" json:"elevation"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// providers
type ProviderWithResource struct {
	ID         string    `db:"id" json:"id"`
	Owner      string    `db:"owner" json:"owner"`
	AreaID     string    `db:"area_id" json:"area_id"`
	RemoteAddr string    `db:"remote_addr" json:"remote_addr"`
	Ip         string    `db:"ip" json:"ip"`
	State      int32     `db:"state" json:"state"`
	Cpu        string    `db:"cpu" json:"cpu"`
	Gpu        string    `db:"gpu" json:"gpu"`
	Memory     string    `db:"memory" json:"memory"`
	Storage    string    `db:"storage" json:"storage"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type User struct {
	ID              int64     `db:"id" json:"id"`
	Uuid            string    `db:"uuid" json:"uuid"`
	Avatar          string    `db:"avatar" json:"avatar"`
	Username        string    `db:"username" json:"username"`
	PassHash        string    `db:"pass_hash" json:"pass_hash"`
	UserEmail       string    `db:"user_email" json:"user_email"`
	WalletAddress   string    `db:"wallet_address" json:"wallet_address"`
	Role            int32     `db:"role" json:"role"`
	AllocateStorage int32     `db:"allocate_storage" json:"allocate_storage"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt       time.Time `db:"deleted_at" json:"deleted_at"`
	ProjectID       int32     `db:"project_id" json:"project_id"`
	ReferralCode    string    `db:"referral_code" json:"referral_code"`
	Referrer        string    `db:"referrer" json:"referrer"`
	ReferrerUserID  string    `db:"referrer_user_id" json:"referrer_user_id"`
}
