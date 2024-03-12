package main

import (
	"bytes"
	"log/slog"

	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/infra"
)

func main() {
	var buf bytes.Buffer
	creator := slogjson.Creator()
	log := slog.New(creator.NewHandler(&buf, infra.SimpleOptions()))
	log.Info("message")
}
