package file

// STDIO is a special constant representing standard input/output streams.
// When used as a filename with ReadFile or WriteFile, it causes the operation
// to use os.Stdin or os.Stdout respectively, instead of opening a file.
//
// This convention is commonly used in Unix command-line tools to allow
// reading from stdin or writing to stdout by specifying "-" as the filename.
const (
	STDIO = "-"
)
