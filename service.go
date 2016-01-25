package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/registry"
	"gopkg.in/mgo.v2/bson"
)

func DetectInstalledPrograms() bool {
	// on successful call of this function from request handler
	// one more http request to get installed programs
	var keypath string
	hostinfo, _ := host.HostInfo()
	if hostinfo.OS == "windows" && strings.Contains(string(hostinfo.Platform), "8.1") {
		keypath = `Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`
	} else {
		keypath = `Software\Microsoft\Windows\CurrentVersion\Uninstall`
	}

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keypath, registry.READ)

	if err != nil {
		fmt.Println(err)
	}

	sl, _ := k.ReadSubKeyNames(0)
	sess := GetConnection()
	defer sess.Close()

	collection := sess.DB("sysinfo").C("programs")
	removeErr := collection.DropCollection()

	if removeErr != nil {
		fmt.Println("problem dropping collection" + removeErr.Error())
	}

	for _, val := range sl {
		path := keypath + `\` + val
		subkey, _ := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
		name, _, _ := subkey.GetStringValue("DisplayName")
		version, _, _ := subkey.GetStringValue("DisplayVersion")
		if strings.TrimSpace(name) != "" {
			//fmt.Printf("%s : %s \n", name, version)
			doc := ProgramInfo{}
			doc.ObjectId = bson.NewObjectId()
			doc.Name = name
			if strings.TrimSpace(version) != "" {
				doc.Version = version
			}
			err := collection.Insert(doc)
			if err != nil {
				fmt.Printf("Can't insert document: %v\n", err)
			}
		}
	}
	//fmt.Println(hostinfo)
	return true
}

func DetectHostMachineInfo() bool {
	hostinfo, _ := host.HostInfo()
	fmt.Println(hostinfo)
	vmem, _ := mem.VirtualMemory()
	fmt.Println(vmem)
	cpuinfo, _ := cpu.CPUInfo()
	fmt.Println(cpuinfo)
	return true
}

func GetProcessInfo() bool {
	// use wmic process command
	r := regexp.MustCompile("[^\\s\\s]+")
	name := "wmic.exe"
	arg0 := "process"
	arg1 := "get"
	arg2 := `description,`
	arg3 := "executablepath"
	cmd := exec.Command(name, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("error executing wmic command" + err.Error())
	}
	op := string(stdout)
	lines := strings.Split(op, "\n")
	sess := GetConnection()
	defer sess.Close()

	collection := sess.DB("sysinfo").C("processes")
	removeErr := collection.DropCollection()
	if removeErr != nil {
		fmt.Println("problem dropping collection" + removeErr.Error())
	}

	for idx, line := range lines {
		if idx == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		details := r.FindAllString(line, -1)
		lenth := len(details)
		if lenth > 0 {
			//fmt.Println(details[0])
			doc := ProcessInfo{}
			doc.ObjectId = bson.NewObjectId()
			doc.Description = details[0]
			if lenth > 1 {
				doc.ExecutablePath = strings.Join(details[1:], "")
			} else {
				doc.ExecutablePath = "Not Available"
			}

			err := collection.Insert(doc)
			if err != nil {
				fmt.Printf("Can't insert document: %v\n", err)
			}
		}
	}
	return true
}
