package bin

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetConfigData() map[string]string {
	SaveLog(log.Fields{
		"group": "server",
	}, log.TraceLevel, "Collecting configuration data...")
	data := make(map[string]string)
	var key, value string
	file, err := os.Open("data/config.txt")
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		CheckErr(scanner.Err())
		subdata := scanner.Text()
		for index, char := range subdata {
			if string(char) == ":" {
				value = subdata[index+2:]
				break
			} else {
				key += string(char)
			}
		}
		data[key] = value
		key = ""
	}
	SaveLog(log.Fields{
		"group": "server",
	}, log.InfoLevel, "Configuration data collected")
	return data
}
