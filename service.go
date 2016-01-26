package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
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
			if name == "Google Chrome" {
				// bad way
				latest := getLatestBrowserVersion("Chrome")
				doc.LatestVersion = latest
				re := regexp.MustCompile("[0-9]+")
				insver := re.FindAllString(version, -1)
				chrom_ver := strings.Join(insver, "")[0:2]
				f1, _ := strconv.Atoi(chrom_ver)
				latver := re.FindAllString(latest, -1)
				f2, _ := strconv.Atoi(strings.Join(latver, ""))
				doc.NeedsUpdate = f2 > f1
			}
			err := collection.Insert(doc)
			if err != nil {
				fmt.Printf("Can't insert document: %v\n", err)
			}
		}
	}

	javainfo := getJavaDetails()
	err = collection.Insert(javainfo)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
	}
	ieinfo := getIEDetails()
	err = collection.Insert(ieinfo)
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
	}
	return true
}

func DetectHostMachineInfo() bool {

	hostinfo, _ := host.HostInfo()

	vmem, _ := mem.VirtualMemory()

	cpuStat := cpu.CPUInfoStat{}
	cpuinfo, _ := cpu.CPUInfo()

	cpuStat.Cores = cpuinfo[0].Cores
	cpuStat.CPU = cpuinfo[0].CPU
	cpuStat.Family = cpuinfo[0].Family
	cpuStat.Model = cpuinfo[0].Model
	cpuStat.ModelName = cpuinfo[0].ModelName
	cpuStat.Mhz = cpuinfo[0].Mhz

	// insert
	sess := GetConnection()
	defer sess.Close()

	collection := sess.DB("sysinfo").C("machineinfo")
	removeErr := collection.DropCollection()

	if removeErr != nil {
		fmt.Println("problem dropping collection" + removeErr.Error())
	}
	err := collection.Insert(cpuStat)
	if err != nil {
		fmt.Printf("Can't insert cpu stat document: %v\n", err)
	}
	err = collection.Insert(vmem)
	if err != nil {
		fmt.Printf("Can't insert mem stat document: %v\n", err)
	}
	err = collection.Insert(hostinfo)
	if err != nil {
		fmt.Printf("Can't insert host stat document: %v\n", err)
	}
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

func getJavaDetails() *ProgramInfo {
	cmdname := "java"
	arg0 := "-version"
	cmd := exec.Command(cmdname, arg0)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error executing wmic command" + err.Error())
	}
	op := string(stdout)
	ver := strings.Split(op, "\n")[0]
	ver = strings.Split(ver, " ")[2]
	installed := string(ver[1 : len(ver)-1])
	doc := ProgramInfo{}
	doc.Name = "Java Runtime"
	doc.ObjectId = bson.NewObjectId()
	doc.Version = installed

	resp, err := http.Get("http://java.com/applet/JreCurrentVersion2.txt")
	if err != nil {
		fmt.Printf("Error on http Get: %v\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	latest := strings.TrimSpace(string(body))
	doc.LatestVersion = latest
	re := regexp.MustCompile("[0-9]+")
	insver := re.FindAllString(installed, -1)
	f1, _ := strconv.Atoi(strings.Join(insver, ""))
	latver := re.FindAllString(latest, -1)
	f2, _ := strconv.Atoi(strings.Join(latver, ""))
	doc.NeedsUpdate = f2 > f1
	return &doc
}

func getIEDetails() *ProgramInfo {
	keypath := `Software\Microsoft\Internet Explorer`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keypath, registry.READ)

	if err != nil {
		fmt.Println(err)
	}
	ieversion, _, _ := k.GetStringValue("Version")
	latestIeVer, _, _ := k.GetStringValue("svcVersion")
	doc := ProgramInfo{}
	doc.Name = "Internet Explorer"
	doc.ObjectId = bson.NewObjectId()
	doc.Version = ieversion
	a, _ := strconv.Atoi(strings.Split(ieversion, `.`)[0])
	b, _ := strconv.Atoi(strings.Split(latestIeVer, `.`)[0])
	if a > b {
		doc.NeedsUpdate = false
	} else if a < b {
		doc.NeedsUpdate = true
	} else {
		c, _ := strconv.Atoi(strings.Split(ieversion, `.`)[1])
		d, _ := strconv.Atoi(strings.Split(latestIeVer, `.`)[1])
		if c >= d {
			doc.NeedsUpdate = false
		} else {
			doc.NeedsUpdate = true
		}
	}
	doc.LatestVersion = latestIeVer
	return &doc
}

func getLatestBrowserVersion(browsername string) string {
	resp, err := http.Get("http://www.webvakman.nl/api/recentbrowserversions")
	if err != nil {
		fmt.Printf("Error on http Get: %v\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(reflect.TypeOf(body))
	var i map[string]interface{}
	e := json.Unmarshal(body, &i)
	if e != nil {
		fmt.Println(e.Error())
	}
	m1 := i["BrowserList"]
	//fmt.Println(reflect.TypeOf(m1))
	m2 := m1.(map[string]interface{})["Windows"].(map[string]interface{})[browsername]
	versn := m2.(map[string]interface{})["LatestVersion"]
	v := versn.(float64)
	var intver int = int(v)
	stringver := strconv.Itoa(intver)
	return stringver
}
