package koyeb

import (
	"encoding/json"
	"time"
)

const RuntimeInfoPath = "/runtime"
const RuntimeInfoPort = 6602

type RuntimeInfo struct {
	GeneratedAt time.Time `json:"generated_at"`
}

func GenerateRuntimeInfoPayloadJSON() ([]byte, error) {
	return json.Marshal(RuntimeInfo{
		GeneratedAt: time.Now(),
	})
}
