package inverted_index

type IndexBuildSourceType string

const (
	FileSourceType   = "file"
	StringSourceType = "string" // for in-memory testings
)

type StringSource struct {
	Name string
	Text string
}
