package model

type AgentDataInterface struct {
	Time           string
	Op             string
	FullPath       string
	FileKey        string
	AppId          string
	Bucket         string
	Retcode        string
	Retmsg         string
	Step           int
	StepIp         string
	Flow           string
	FileSize       int64
	RequestId      string
	ClientIp       string
	CosTime        int
	DownloadLength int
	FileCtime      int
	ContentType    string
	StorageType    string
	HttpStatusCode int
	RequestHttp    string
}
