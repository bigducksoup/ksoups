package message

type DataType string

const (
	DEFAULT DataType = "DEFAULT"
	ERROR   DataType = "ERROR"
)

const (
	READDIR     DataType = "DIRREAD"
	READDIRRESP DataType = "DIRREADRESP"
	CREATE_DIR  DataType = "CREATE_DIR"
)

const (
	READFILE       DataType = "READFILE"
	READFILERESP   DataType = "READFILERESP"
	MODIFYFILE     DataType = "FILEMODIFY"
	MODIFYFILERESP DataType = "FILEMODIFYRESP"
	CREATEFILE     DataType = "CREATEFILE"
	DELETEFILE     DataType = "DELETEFILE"
)

const (
	RUN_SHORTCUT DataType = "RUN_SHORTCUT"
)
