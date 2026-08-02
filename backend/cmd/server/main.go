// Package main is the Wire composition root for the VEDO EduTrack backend.
//
// It wires together all 10 bounded contexts (Clean Architecture modules)
// and platform adapters into a single modular monolith binary.
//
// See ADR-DES.INFRA.modular-monolith-approach and
// ADR-IMPL.PROCESS.repository-structure for architecture decisions.
package main

import (
	"context"
	"log"
	"os"
)

func main() {
	log.Println("[INFO] [server.main] starting VEDO EduTrack backend")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Wire all modules and platform adapters via google/wire.
	// cfg, err := platformconfig.Load(ctx)
	// lgr := platformlogger.New(cfg.LogLevel)
	// db, err := platformpostgres.Connect(ctx, cfg.DatabaseURL)
	// ...wire each module...

	log.Printf("[INFO] [server.main] VEDO EduTrack backend running on PID %d", os.Getpid())

	// Block until signal.
	<-ctx.Done()
	log.Println("[INFO] [server.main] VEDO EduTrack backend shutting down")
}
