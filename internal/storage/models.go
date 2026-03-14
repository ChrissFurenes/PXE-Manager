package storage

type Asset struct {
	ID          int64
	Name        string
	FileName    string
	FilePath    string
	FileType    string
	SizeBytes   int64
	Description string
}

type Profile struct {
	ID        int64
	Name      string
	BootMode  string
	BootType  string
	Kernel    string
	Initrd    string
	ImagePath string
	Cmdline   string
	Enabled   bool
}

type Client struct {
	ID          int64
	MAC         string
	Hostname    string
	ProfileID   int64
	ShowMenu    bool
	Description string
}

type Setting struct {
	Key   string
	Value string
}
