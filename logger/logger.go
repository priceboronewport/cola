package logger

import (
	"github.com/priceboronewport/cola/filestore"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var f *os.File
var l *log.Logger
var config *filestore.FileStore

func Init(config_filename string) {
	config = filestore.New(config_filename)
	log_filename := config.Read("log_filename")
	if log_filename != "" {
		var err error
		f, err = os.OpenFile(log_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			l = log.New(f, "", log.LstdFlags)
		} else {
			log.Println(err)
		}
	}
}

func Emergency(strs ...string) {
	Output(0, strs)
}

func Alert(strs ...string) {
	Output(1, strs)
}

func Critical(strs ...string) {
	Output(2, strs)
}

func Error(strs ...string) {
	Output(3, strs)
}

func Warning(strs ...string) {
	Output(4, strs)
}

func Notice(strs ...string) {
	Output(5, strs)
}

func Informational(strs ...string) {
	Output(6, strs)
}

func Debug(strs ...string) {
	Output(7, strs)
}

func Output(output_level int, strs []string) {
	log_level := 6
	if config != nil {
		log_level, _ = strconv.Atoi(config.Read("log_level"))
	}
	if log_level >= output_level {
		output := fmt.Sprintf("%d", output_level)
		for _, str := range strs {
			output += " - " + str
		}
		if l != nil {
			l.Println(output)
		} else {
			log.Println(output)
		}
		if output_level == 0 {
			panic(errors.New(output))
		}
	}
}
