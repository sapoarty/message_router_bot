package config

import (
  "github.com/joho/godotenv"
  "message_router_bot/structures"
  "log"
)

var UserStates map[int]structures.User

func Init() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}
