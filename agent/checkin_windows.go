package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

func (d *agentDaemon) checkin() {

	var data checkinData

	data.ID = d.ID
	data.Version = d.version

	sd := d.getSystemData()
	if sd != nil {
		data.Serial = "TODO"
	}

	b := &bytes.Buffer{}
	ge := gob.NewEncoder(b)
	err := ge.Encode(data)
	if checkError(err) {
		return
	}

	resp, err := d.hc.Post(fmt.Sprintf("%s/%s", d.cmdHost, checkinURL), "application/octet-stream", b)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

}

func (d *agentDaemon) getSystemData() *AppleSystemProfilerOutput {

	v, ok := d.systemData.(AppleSystemProfilerOutput)
	if ok {
		return &v
	} else {
		return nil
	}
}

func (d *agentDaemon) sendSystemData() {
	var data checkinData
	data.ID = d.ID

	output, err := run("hostname")
	if checkError(err) {
		return
	}

	d.hostname = output

	data.Hostname = output
	data.OS = runtime.GOOS

	// start := time.Now()
	log.Printf("Reading system data...")
	d.systemData, err = readSystemData()
	if checkError(err) {
		return
	}
	data.Payload = d.systemData
	data.Version = d.version
	sd := d.getSystemData()
	if sd == nil {
		log.Printf("System data is nil")
		return
	}
	data.Serial = "TODO"

	// log.Printf("Got system data (took %s): %+v", time.Since(start).String(), d.getSystemData())
	//log.Printf("System data: %+v", d.getSystemData())

	b := &bytes.Buffer{}
	gob.Register(data)
	gob.Register(data.Payload)
	ge := gob.NewEncoder(b)
	err = ge.Encode(data)
	if checkError(err) {
		return
	}

	resp, err := d.hc.Post(fmt.Sprintf("%s/%s", d.cmdHost, systemDataURL), "application/octet-stream", b)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

	var cr checkinResponse

	gob.Register(cr)
	gd := gob.NewDecoder(resp.Body)

	err = gd.Decode(&cr)
	if checkError(err) {
		return
	}

	d.ID = cr.ID

	for _, c := range cr.Commands {
		log.Printf("Received command named '%s' with arguments: %#v", c.Name, c.Arguments)
	}

}

type Command struct {
	Name      string
	Arguments []string
}

type checkinResponse struct {
	ID       int
	Commands []Command
}

func readSystemData() (*computerInfo, error) {
    

	jsonData, err := run(fmt.Sprintf("get-computerinfo | convertto-json"))
	if checkError(err) {
		return nil, err
	}

    wci := &windowsComputerInfo{}

    err = json.Unmarshal(jsonData, wci)
    if checkError(err){
        return nil, err
    }

    return ci, nil
}

type windowsComputerInfo struct {

}

// Actually, use get-computerinfo | convertto-json

// use msinfo32 /nfo C:\Windows\temp\output.xml
// or use powershell: get-ciminstance -class win32_operatingsystem | convertto-json > C:\Windows\temp\output.json
