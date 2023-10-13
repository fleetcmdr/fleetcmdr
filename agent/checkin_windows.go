package main

func (d *daemon) checkin() {

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
