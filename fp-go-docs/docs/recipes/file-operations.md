---
sidebar_position: 11
title: File Operations
description: Functional file I/O patterns
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 11 / 17"
  title="File"
  titleAccent="Operations"
  lede="Perform file I/O operations using functional patterns with IOEither for lazy evaluation, resource safety, and composable file operations."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '7' },
    { label: 'Use Cases', value: 'File I/O, Data Processing, Logs' }
  ]}
/>

<TLDR>
  <TLDRCard title="Lazy Evaluation" icon="clock">
    Operations don't execute until called—compose file operations without triggering side effects.
  </TLDRCard>
  <TLDRCard title="Resource Safety" icon="shield">
    Use bracket for proper cleanup—ensures files are closed even when errors occur.
  </TLDRCard>
  <TLDRCard title="Composability" icon="layers">
    Chain multiple file operations—build complex workflows from simple operations.
  </TLDRCard>
</TLDR>

<Section id="reading-files" number="01" title="Reading" titleAccent="Files">

Read files with proper error handling and resource management.

<CodeCard file="read-file.go">
{`package main

import (
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func readFile(path string) IOE.IOEither[error, []byte] {
    return IOE.TryCatch(func() ([]byte, error) {
        return os.ReadFile(path)
    })
}

func readFileAsString(path string) IOE.IOEither[error, string] {
    return IOE.Map(func(data []byte) string {
        return string(data)
    })(readFile(path))
}

func main() {
    result := readFileAsString("config.json")()
    
    if result.IsLeft() {
        fmt.Println("Error reading file:", result.Left())
    } else {
        fmt.Println("File contents:", result.Right())
    }
}`}
</CodeCard>

<CodeCard file="read-lines.go">
{`package main

import (
    "bufio"
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func readLines(path string) IOE.IOEither[error, []string] {
    return IOE.TryCatch(func() ([]string, error) {
        file, err := os.Open(path)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        
        var lines []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            lines = append(lines, scanner.Text())
        }
        
        if err := scanner.Err(); err != nil {
            return nil, err
        }
        
        return lines, nil
    })
}

func countLines(path string) IOE.IOEither[error, int] {
    return F.Pipe1(
        readLines(path),
        IOE.Map(func(lines []string) int {
            return len(lines)
        }),
    )
}

func main() {
    result := countLines("data.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("File has %d lines\\n", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="writing-files" number="02" title="Writing" titleAccent="Files">

Write and append to files with proper error handling.

<CodeCard file="write-file.go">
{`package main

import (
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func writeFile(path string, data []byte) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        err := os.WriteFile(path, data, 0644)
        if err != nil {
            return 0, err
        }
        return len(data), nil
    })
}

func writeString(path, content string) IOE.IOEither[error, int] {
    return writeFile(path, []byte(content))
}

func main() {
    content := "Hello, functional world!"
    result := writeString("output.txt", content)()
    
    if result.IsLeft() {
        fmt.Println("Error writing file:", result.Left())
    } else {
        fmt.Printf("Wrote %d bytes\\n", result.Right())
    }
}`}
</CodeCard>

<CodeCard file="append-file.go">
{`package main

import (
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func appendToFile(path, content string) IOE.IOEither[error, int] {
    return IOE.TryCatch(func() (int, error) {
        file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return 0, err
        }
        defer file.Close()
        
        return file.WriteString(content)
    })
}

func appendLine(path, line string) IOE.IOEither[error, int] {
    return appendToFile(path, line+"\\n")
}

func main() {
    result := appendLine("log.txt", "New log entry")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Appended %d bytes\\n", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="resource-management" number="03" title="Resource" titleAccent="Management">

Use bracket for safe resource handling that guarantees cleanup.

<CodeCard file="bracket-pattern.go">
{`package main

import (
    "bufio"
    "fmt"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func withFile[A any](
    path string,
    use func(*os.File) IOE.IOEither[error, A],
) IOE.IOEither[error, A] {
    return IOE.Bracket(
        // Acquire resource
        IOE.TryCatch(func() (*os.File, error) {
            return os.Open(path)
        }),
        // Use resource
        use,
        // Release resource (always called)
        func(file *os.File) IOE.IOEither[error, struct{}] {
            return IOE.TryCatch(func() (struct{}, error) {
                return struct{}{}, file.Close()
            })
        },
    )
}

func readFirstLine(path string) IOE.IOEither[error, string] {
    return withFile(path, func(file *os.File) IOE.IOEither[error, string] {
        return IOE.TryCatch(func() (string, error) {
            scanner := bufio.NewScanner(file)
            if scanner.Scan() {
                return scanner.Text(), nil
            }
            if err := scanner.Err(); err != nil {
                return "", err
            }
            return "", fmt.Errorf("file is empty")
        })
    })
}

func main() {
    result := readFirstLine("data.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("First line:", result.Right())
    }
}`}
</CodeCard>

<CodeCard file="copy-file.go">
{`package main

import (
    "fmt"
    "io"
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func copyFile(src, dst string) IOE.IOEither[error, int64] {
    return IOE.Bracket(
        // Open source file
        IOE.TryCatch(func() (*os.File, error) {
            return os.Open(src)
        }),
        // Copy to destination
        func(srcFile *os.File) IOE.IOEither[error, int64] {
            return IOE.Bracket(
                // Create destination file
                IOE.TryCatch(func() (*os.File, error) {
                    return os.Create(dst)
                }),
                // Copy data
                func(dstFile *os.File) IOE.IOEither[error, int64] {
                    return IOE.TryCatch(func() (int64, error) {
                        return io.Copy(dstFile, srcFile)
                    })
                },
                // Close destination
                func(dstFile *os.File) IOE.IOEither[error, struct{}] {
                    return IOE.TryCatch(func() (struct{}, error) {
                        return struct{}{}, dstFile.Close()
                    })
                },
            )
        },
        // Close source
        func(srcFile *os.File) IOE.IOEither[error, struct{}] {
            return IOE.TryCatch(func() (struct{}, error) {
                return struct{}{}, srcFile.Close()
            })
        },
    )
}

func main() {
    result := copyFile("source.txt", "destination.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error copying file:", result.Left())
    } else {
        fmt.Printf("Copied %d bytes\\n", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="directory-operations" number="04" title="Directory" titleAccent="Operations">

List and walk directory trees with functional patterns.

<CodeCard file="list-directory.go">
{`package main

import (
    "fmt"
    "os"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
)

func listDir(path string) IOE.IOEither[error, []os.FileInfo] {
    return IOE.TryCatch(func() ([]os.FileInfo, error) {
        file, err := os.Open(path)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        
        return file.Readdir(-1)
    })
}

func listFiles(path string) IOE.IOEither[error, []string] {
    return F.Pipe2(
        listDir(path),
        IOE.Map(func(entries []os.FileInfo) []string {
            return A.FilterMap(func(entry os.FileInfo) O.Option[string] {
                if !entry.IsDir() {
                    return O.Some(entry.Name())
                }
                return O.None[string]()
            })(entries)
        }),
    )
}

func main() {
    result := listFiles(".")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        files := result.Right()
        fmt.Printf("Found %d files:\\n", len(files))
        for _, file := range files {
            fmt.Println(" -", file)
        }
    }
}`}
</CodeCard>

</Section>

<Section id="transformations" number="05" title="File" titleAccent="Transformations">

Transform file contents with functional pipelines.

<CodeCard file="transform-file.go">
{`package main

import (
    "fmt"
    "strings"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func transformFile(
    input, output string,
    transform func(string) string,
) IOE.IOEither[error, int] {
    return F.Pipe3(
        readFileAsString(input),
        IOE.Map(transform),
        IOE.Chain(func(content string) IOE.IOEither[error, int] {
            return writeString(output, content)
        }),
    )
}

func toUpperCase(input, output string) IOE.IOEither[error, int] {
    return transformFile(input, output, strings.ToUpper)
}

func replaceInFile(input, output, old, new string) IOE.IOEither[error, int] {
    return transformFile(input, output, func(content string) string {
        return strings.ReplaceAll(content, old, new)
    })
}

func main() {
    result := toUpperCase("input.txt", "output.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Transformed file, wrote %d bytes\\n", result.Right())
    }
}`}
</CodeCard>

<CodeCard file="process-lines.go">
{`package main

import (
    "fmt"
    "strings"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func processLines(
    input, output string,
    process func(string) string,
) IOE.IOEither[error, int] {
    return F.Pipe3(
        readLines(input),
        IOE.Map(func(lines []string) []string {
            return A.Map(process)(lines)
        }),
        IOE.Chain(func(lines []string) IOE.IOEither[error, int] {
            content := strings.Join(lines, "\\n")
            return writeString(output, content)
        }),
    )
}

func removeEmptyLines(input, output string) IOE.IOEither[error, int] {
    return F.Pipe3(
        readLines(input),
        IOE.Map(func(lines []string) []string {
            return A.Filter(func(line string) bool {
                return strings.TrimSpace(line) != ""
            })(lines)
        }),
        IOE.Chain(func(lines []string) IOE.IOEither[error, int] {
            content := strings.Join(lines, "\\n")
            return writeString(output, content)
        }),
    )
}

func main() {
    result := removeEmptyLines("input.txt", "output.txt")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Processed file, wrote %d bytes\\n", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="batch-operations" number="06" title="Batch" titleAccent="Operations">

Process multiple files sequentially or in parallel.

<CodeCard file="batch-processing.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func processFiles(
    files []string,
    process func(string) IOE.IOEither[error, int],
) IOE.IOEither[error, []int] {
    return A.Traverse[string](IOE.Applicative[error, int]())(
        process,
    )(files)
}

func countLinesInFiles(files []string) IOE.IOEither[error, int] {
    return F.Pipe2(
        processFiles(files, func(file string) IOE.IOEither[error, int] {
            return countLines(file)
        }),
        IOE.Map(func(counts []int) int {
            return A.Reduce(func(acc, n int) int {
                return acc + n
            })(0)(counts)
        }),
    )
}

func main() {
    files := []string{"file1.txt", "file2.txt", "file3.txt"}
    result := countLinesInFiles(files)()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Printf("Total lines across all files: %d\\n", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="temporary-files" number="07" title="Temporary" titleAccent="Files">

Create and use temporary files with automatic cleanup.

<CodeCard file="temp-files.go">
{`package main

import (
    "fmt"
    "os"
    "strings"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func withTempFile[A any](
    pattern string,
    use func(*os.File) IOE.IOEither[error, A],
) IOE.IOEither[error, A] {
    return IOE.Bracket(
        // Create temp file
        IOE.TryCatch(func() (*os.File, error) {
            return os.CreateTemp("", pattern)
        }),
        // Use temp file
        use,
        // Clean up
        func(file *os.File) IOE.IOEither[error, struct{}] {
            return IOE.TryCatch(func() (struct{}, error) {
                name := file.Name()
                file.Close()
                return struct{}{}, os.Remove(name)
            })
        },
    )
}

func processWithTempFile(data string) IOE.IOEither[error, string] {
    return withTempFile("process-*.txt", func(temp *os.File) IOE.IOEither[error, string] {
        return F.Pipe3(
            // Write to temp file
            IOE.TryCatch(func() (int, error) {
                return temp.WriteString(data)
            }),
            // Process temp file
            IOE.Chain(func(_ int) IOE.IOEither[error, string] {
                return readFileAsString(temp.Name())
            }),
            // Transform result
            IOE.Map(strings.ToUpper),
        )
    })
}

func main() {
    result := processWithTempFile("hello world")()
    
    if result.IsLeft() {
        fmt.Println("Error:", result.Left())
    } else {
        fmt.Println("Result:", result.Right())
    }
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="08" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Always use bracket for resources** — Ensures cleanup even when errors occur
  </ChecklistItem>
  <ChecklistItem status="required">
    **Check file existence** — Verify files exist before operations
  </ChecklistItem>
  <ChecklistItem status="required">
    **Use appropriate permissions** — Set correct file permissions (0644, 0600, etc.)
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Handle large files carefully** — Use streaming for large files to avoid memory issues
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Validate paths** — Check and sanitize file paths before use
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Use temp files for safety** — Process data in temp files before overwriting originals
  </ChecklistItem>
</Checklist>

</Section>
