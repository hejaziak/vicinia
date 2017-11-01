package main

import(
    "github.com/satori/go.uuid"
)

type WelcomeStruct struct {
    Message      string    `json:"message"`
    Uuid uuid.UUID      `json:"uuid"`
}