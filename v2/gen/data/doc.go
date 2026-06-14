package data

//go:generate go run ../main.go cp --src ../../../skills --dst ./skills

//go:generate go run ../main.go embed --package data --map-name Skills --src ./skills --dst ./gen_skills.go
