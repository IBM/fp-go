# Skills Directory

This directory contains skill documentation files that are **redundantly stored** for embedding purposes.

## Important Note

The content in this directory is **NOT the normative source**. These files are copies maintained for Go embedding via the `//go:embed` directive.

## Normative Source

The **authoritative source** for these skill files is located in the main fp-go repository:

**[https://github.com/IBM/fp-go/tree/main/skills](https://github.com/IBM/fp-go/tree/main/skills)**

All skill documentation files are maintained in the central skills directory of the fp-go repository.

## Purpose

These files are embedded into the binary at compile time through [`data/gen_skills.go`](../gen_skills.go), which uses Go's `//go:embed` directive to include the skill documentation as byte arrays. This allows the application to access skill documentation without external file dependencies at runtime.

## Updating Skills

When updating skill documentation:

1. Update the normative source at [https://github.com/IBM/fp-go/tree/main/skills](https://github.com/IBM/fp-go/tree/main/skills)
2. Copy the updated content to this directory
3. Run `go generate ./...` from the `v2/gen` project root to regenerate the embedding code and keep sources in sync

Do not edit these files directly unless you intend to update the normative source as well.