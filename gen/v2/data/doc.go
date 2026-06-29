package data

//go:generate go run ../main.go cp --src ../../../skills --dst ./skills

//go:generate go run ../main.go examples ingest --verbose --src ../../../v2 --db ./examples/examples.db

//go:generate go run ../main.go embed --package data --map-name Skills --src ./skills --dst ./gen_skills.go

//go:generate go run ../main.go embed --package data --map-name Examples --src ./examples --dst ./gen_examples.go
