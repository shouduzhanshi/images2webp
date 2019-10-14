package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var utilsPath string

//"-q", "80", "-mt", "-v", "-progress", "-o"
func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cmdParms := os.Args
	utilsPath = dir + "/cwebp"
	projectPath := cmdParms[1]
	isReplace, _ := strconv.ParseBool(cmdParms[2])
	parms := make([]string, 0)
	for i := 3; i < len(cmdParms); i++ {
		parms = append(parms, cmdParms[i])
	}
	parms = append(parms, "-o")
	readDir(projectPath, isReplace, parms)
}

func readDir(path string, isReplace bool, parms []string) {
	if infos, e := ioutil.ReadDir(path); e == nil {
		for _, temp := range infos {
			dir := temp.IsDir()
			name := temp.Name()
			if dir && !strings.Contains(name, "assets") {
				readDir(path+"/"+name, false, parms)
			} else if strings.HasSuffix(name, ".JPG") || strings.HasSuffix(name, ".JPEG") ||
				strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") ||
				strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".PNG") {
				split := strings.Split(name, ".")
				dispFile(path+"/"+name, path, split[0], parms)
				if !isReplace {
					os.Remove(path + "/" + name)
				}
			}
		}
	}
}

func dispFile(path, dirPath, name string, parms []string) {
	parms = append(parms, dirPath+"/"+name+".webp")
	parms = append(parms, path)
	execCmd(utilsPath, parms)
}

func execCmd(shell string, raw []string) (int, error) {
	fmt.Println(shell, raw)
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return 0, nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for s.Scan() {
		text := s.Text()
		fmt.Println(text)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
	return 0, nil
}
