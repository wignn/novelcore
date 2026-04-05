package main

import "time"

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	AvatarUrl *string   `json:"avatarUrl,omitempty"`
	Bio       *string   `json:"bio,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
