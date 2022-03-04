package logger

import (
	"log"
	"os"
	"sync"
)

type logger struct {
	fileName string
	*log.Logger
}

var once sync.Once
var Logger = setLogger()

func setLogger() *logger {
	var logger *logger
	once.Do(func() {
		file, err := os.OpenFile("./log/err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

		if err != nil {
			panic(err)
		}

		log.SetOutput(file)

		logger.fileName = file.Name()
		logger.Logger = log.New(file, "school helper: ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return logger
}
