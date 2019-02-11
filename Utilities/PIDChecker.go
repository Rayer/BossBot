package Utilities

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
)

type ExecutionBlock func()

func ExecuteCode(fullpath string, mainLoop ExecutionBlock) {

	fileStat, err := os.Stat(fullpath)
	if err == nil && fileStat.IsDir() {
		fullpath = fullpath + "/" + os.Args[0] + ".pid"
	}

	log.Debugln("Trying to handle pid file : " + fullpath)
	if _, err := os.Stat(fullpath); err == nil {

		//Read PID
		fileContent, err := ioutil.ReadFile(fullpath)
		if err != nil {
			panic("PID file exist, but corrupted")
		}
		pid, err := strconv.Atoi(string(fileContent))
		if err != nil {
			log.Errorln("PID file corrupted, delete it")
		}

		proc, _ := os.FindProcess(pid)
		if proc != nil {
			err = proc.Kill()
			if err != nil {
				log.Warnln("PID file exist but fail to kill process, is there a permission issue?")
			} else {
				log.Infof("Killing original process(%d) and replacing this one...", pid)
			}
		} else {
			log.Warnln("PID file exists but relative PID not exists, going on....")
		}

		_ = os.Remove(fullpath)
	}

	pid := os.Getpid()
	err = ioutil.WriteFile(fullpath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Warnln("Failed to write PIDFile : " + fullpath)
	}

	mainLoop()

}
