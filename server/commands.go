package main

import (
	"context"
	"database/sql"
	"html/template"
	"time"
)

func (d *serverDaemon) getAgentCommands(id int) (commands []Command, err error) {
	q := "SELECT id, ts, input, COALESCE(output,''), scheduled_ts, delivered_ts, executed_ts FROM commands WHERE agent_id = $1 ORDER BY scheduled_ts DESC LIMIT 20"
	rows, err := d.db.QueryContext(context.Background(), q, id)
	if checkError(err) {
		return
	}

	cmds := []Command{}

	for rows.Next() {
		c := Command{}
		err = rows.Scan(&c.ID, &c.TS, &c.Input, &c.Output, &c.ScheduledTS, &c.DeliveredTS, &c.ExecutedTS)
		if checkError(err) {
			return
		}

		if c.ExecutedTS.Valid {
			c.Executed = true
		}

		cmds = append([]Command{c}, cmds...)
	}

	return cmds, nil
}

type Command struct {
	ID          int
	UUID        string
	TS          time.Time
	Input       string
	Output      string
	ScheduledTS sql.NullTime
	DeliveredTS sql.NullTime
	ExecutedTS  sql.NullTime
	Executed    bool
	Special     int64
}

const (
	SpecialUpgrade = 1 << iota
)

type darwinSystemData struct {
	AgentData               *agent
	SystemData              AppleSystemProfilerOutput
	Commands                []Command
	Scripts                 []Script
	CheckinHistorySparkline template.HTML
}
