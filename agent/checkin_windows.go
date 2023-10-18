package main

func (d *agentDaemon) checkin() {

	var data checkinData

	getSerial()

	b := &bytes.Buffer{}
	ge := gob.NewEncoder(b)
	ge.Encode()

	resp, err := d.hc.Post(d.cmdr, "application/octet-stream", b)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

}

// use msinfo32 /nfo C:\Windows\temp\output.xml
// or use powershell: get-ciminstance -class win32_operatingsystem | convertto-json > C:\Windows\temp\output.json
