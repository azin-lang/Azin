package main

import (
	"log"

	"github.com/azin-lang/Azin/internal/codegen/x86_64"
	"github.com/azin-lang/Azin/internal/object/elf"
)

func main() {
	code := x86_64.EmitMinimalExit()

	if err := elf.WriteRelocatableObject("output", code.Emitter); err != nil {
		log.Fatalf("failed to emit executable: %v", err)
	}
}
