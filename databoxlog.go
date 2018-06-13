package libDatabox

import (
	"encoding/json"
	"log"
	"runtime"
	"strconv"
)

type Logger struct {
	Store *CoreStoreClient
}

type LogEntries struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type Logs []LogEntries

var debug = false

func New(store *CoreStoreClient, outputDebugLogs bool) (*Logger, error) {

	dsmd := DataSourceMetadata{
		Description:    "container manager logs",
		ContentType:    "aplication/json",
		Vendor:         "databox",
		DataSourceType: "databox-logs",
		DataSourceID:   "cmlogs",
		StoreType:      "tsblob",
		IsActuator:     false,
		Unit:           "",
		Location:       "",
	}

	debug = outputDebugLogs

	err := store.RegisterDatasource(dsmd)
	ChkErr(err)

	return &Logger{
		Store: store,
	}, nil
}

func (l Logger) Info(msg string) {
	Info(msg)
	err := l.Store.TSBlobWrite("cmlogs", []byte("{\"log\":"+strconv.Quote(msg)+",\"type\":\"INFO\"}"))
	ChkErr(err)
}
func (l Logger) Warn(msg string) {
	Warn(msg)
	err := l.Store.TSBlobWrite("cmlogs", []byte("{\"log\":"+strconv.Quote(msg)+",\"type\":\"WARN\"}"))
	ChkErr(err)
}
func (l Logger) Err(msg string) {
	Err(msg)
	err := l.Store.TSBlobWrite("cmlogs", []byte("{\"log\":"+strconv.Quote(msg)+",\"type\":\"ERROR\"}"))
	ChkErr(err)
}
func (l Logger) Debug(msg string) {
	Debug(msg)
	err := l.Store.TSBlobWrite("cmlogs", []byte("{\"log\":"+strconv.Quote(msg)+",\"type\":\"DEBUG\"}"))
	ChkErr(err)
}

func (l Logger) ChkErr(err error) {
	if err == nil {
		return
	}
	Err(err.Error())
	l.Err(err.Error())

}

func (l Logger) GetLastNLogEntries(n int) Logs {

	var logs Logs
	data, err := l.Store.TSBlobLastN("cmlogs", n)
	l.ChkErr(err)
	json.Unmarshal(data, &logs)

	return logs
}

func (l Logger) GetLastNLogEntriesRaw(n int) []byte {

	data, err := l.Store.TSBlobLastN("cmlogs", n)
	l.ChkErr(err)
	return data

}

func ChkErr(err error) {
	if err == nil {
		return
	}
	Err(err.Error())
}

func ChkErrFatal(err error) {
	if err == nil {
		return
	}
	log.SetFlags(log.Ldate | log.Ltime)
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	log.Fatal("[ERROR] " + file + " L" + strconv.Itoa(line) + ":" + err.Error())
}

func Info(msg string) {
	log.SetPrefix("[INFO]")
	log.SetFlags(log.LstdFlags)
	log.Println(msg)
}

func Warn(msg string) {
	log.SetPrefix("[WARNING]")
	log.SetFlags(log.LstdFlags)
	log.Println(msg)
}

func Err(msg string) {
	log.SetPrefix("[ERROR]")
	log.SetFlags(log.Ldate | log.Ltime)
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	log.Println(file + " L" + strconv.Itoa(line) + ":" + msg)
}

func Debug(msg string) {
	if debug == true {
		log.SetPrefix("[DEBUG]")
		log.SetFlags(log.LstdFlags)
		log.Println(msg)
	}

}
