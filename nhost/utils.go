package nhost

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func appendEnvVars(payload map[interface{}]interface{}, prefix string) []string {
	var response []string
	for key, item := range payload {
		switch item := item.(type) {
		/*
			case map[interface{}]interface{}:
				response = append(response, appendEnvVars(item, prefix)...)
		*/
		case map[interface{}]interface{}:
			for key, value := range item {
				switch value := value.(type) {
				case map[interface{}]interface{}:
					for newkey, newvalue := range value {
						if newvalue != "" {
							response = append(response, fmt.Sprintf("%s_%v_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), strings.ToUpper(fmt.Sprint(newkey)), newvalue))
						}
					}
				case interface{}, string:
					if value != "" {
						if key.(string) == "smtp_host" {
							response = append(response, fmt.Sprintf("%s_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), GetContainerName(value.(string))))
						} else {
							response = append(response, fmt.Sprintf("%s_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), value))
						}
					}
				}
			}
		case interface{}:
			if item != "" {
				response = append(response, fmt.Sprintf("%s_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), item))
			}
		}
	}
	return response
}

// generate a random 128 byte key
func generateRandomKey() string {
	key := make([]byte, 128)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func GetPort(low, hi int) int {

	// generate a random port value
	port := strconv.Itoa(low + rand.Intn(hi-low))

	// validate wehther the port is available
	if !portAvaiable(port) {
		return GetPort(low, hi)
	}

	// return the value, if it's available
	response, _ := strconv.Atoi(port)
	return response
}

func portAvaiable(port string) bool {

	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return false
	}

	ln.Close()
	return true
}

func GetContainerName(name string) string {
	return strings.Join([]string{PREFIX, name}, "_")
}

func openbrowser(url string) error {

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}