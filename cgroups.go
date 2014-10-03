package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var root = "/sys/fs/cgroup"

func listResources() []string {
	var resources []string

	files, _ := ioutil.ReadDir(root)
	for _, f := range files {
		resources = append(resources, f.Name())
	}
	return resources
}

func cgroupPath(cgroup string) string {
	cgroup = strings.Replace(cgroup, ":", "", -1)
	cgroup = strings.Replace(cgroup, "|", "/", -1)
	return fmt.Sprintf("%s/%s/tasks", root, cgroup)
}

func listPids(cgroup string) []int {
	cmd := exec.Command("cat", cgroupPath(cgroup))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	var pids []int

	for _, pid := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		ipid, _ := strconv.Atoi(pid)
		pids = append(pids, ipid)
	}
	return pids
}

func addPid(cgroup string, pid string) bool {
	echo := exec.Command("echo", pid)
	tee := exec.Command("tee", "-a", cgroupPath(cgroup))
	echoOut, _ := echo.StdoutPipe()
	echo.Start()
	tee.Stdin = echoOut
	_, err := tee.Output()

	if err != nil {
		return false
	}
	return true
}

func recurseSubdirs(top string) []string {
	var children []string
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", root, top))
	for _, f := range files {
		if f.IsDir() {
			current := fmt.Sprintf("%s/%s", top, f.Name())
			children = append(children, current)
			for _, child := range recurseSubdirs(current) {
				children = append(children, child)
			}
		}
	}
	return children
}

func listCgroups() []string {
	var cgroups []string
	re := regexp.MustCompile("/([^/]+)/?(.*)")
	for _, d := range recurseSubdirs("") {
		group := re.ReplaceAllString(d, "$1:/$2")
		cgroups = append(cgroups, strings.Replace(group, "/", "|", -1))
	}
	return cgroups
}
