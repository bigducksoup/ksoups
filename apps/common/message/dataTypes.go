package message

type DataType string

const (
	DEFAULT DataType = "DEFAULT"
	ERROR   DataType = "ERROR"
)

const (
	READDIR         DataType = "DIRREAD"
	READDIRRESP     DataType = "DIRREADRESP"
	CREATE_DIR      DataType = "CREATE_DIR"
	CREATE_DIR_RESP DataType = "CREATE_DIR_RESP"
)

const (
	READ_FILE        DataType = "READFILE"
	READ_FILE_RESP   DataType = "READFILERESP"
	MODIFY_FILE      DataType = "FILEMODIFY"
	MODIFY_FILE_RESP DataType = "FILEMODIFYRESP"
	CREATE_FILE      DataType = "CREATEFILE"
	CREATE_FILE_RESP DataType = "CREATEFILERESP"
	DELETE_FILE      DataType = "DELETEFILE"
	DELETE_FILE_RESP DataType = "DELETEFILERESP"
)

const (
	RUN_SHORTCUT       DataType = "RUN_SHORTCUT"
	RUN_SHORTCUT_RESP  DataType = "RUN_SHORTCUT_RESP"
	CREATE_SCRIPT      DataType = "CREATE_SCRIPT"
	CREATE_SCRIPT_RESP DataType = "CREATE_SCRIPT_RESP"
	SHORTCUT_OUTPUT    DataType = "SHORTCUT_OUTPUT"
)
