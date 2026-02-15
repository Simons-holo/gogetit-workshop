package downloader

const DefaultTimeout = 30

type Config struct {
	OutputDir   string
	Concurrency int
	Timeout     int
	Retry       int
	UserAgent   string
}
