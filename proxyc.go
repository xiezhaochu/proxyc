package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

var configEnv = []string{
	`http_proxy=socks5://127.0.0.1:1080`,
	`https_proxy=socks5://127.0.0.1:1080`,
}

func CurrentHomedir() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func readConfig() {
	homedir, err := CurrentHomedir()
	if err != nil {
		return
	}
	configFile := homedir + `/.proxyc`
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	var configIfc interface{}
	err = json.Unmarshal(configContent, &configIfc)
	if err != nil {
		log.Fatal("config file error: ", err)
	}
	configEnv = []string{}
	mapConfig := configIfc.(map[string]interface{})
	for k, v := range mapConfig {
		configEnv = append(configEnv, k+"="+v.(string))
	}
}

func main() {
	args := make([]string, 0)
	readConfig()
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if len(os.Args) < 2 {
		return
	}
	cmd := exec.Command(os.Args[1], args...)
	cmd.Env = append(os.Environ(), configEnv...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
