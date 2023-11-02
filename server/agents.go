package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type agent struct {
	ID         int
	ClientID   int
	Name       string
	Serial     string
	Deleted    time.Time
	OS         string
	SystemData string
}

func (d *serverDaemon) getAgents(limit, skip int) []*agent {
	q := "SELECT id, client_id, host_name FROM agents WHERE id NOT IN (select id FROM deleted_agents) ORDER BY host_name asc LIMIT $1 OFFSET $2"
	rows, err := d.db.QueryContext(context.Background(), q, limit, skip)
	if checkError(err) {
		return nil
	}

	var agents []*agent
	for rows.Next() {
		a := &agent{}
		err = rows.Scan(&a.ID, &a.ClientID, &a.Name)
		if checkError(err) {
			return nil
		}

		agents = append(agents, a)
	}

	log.Printf("Returning %d agents", len(agents))

	return agents
}

func (d *serverDaemon) getAgentByID(id int) (*agent, error) {
	a := &agent{}
	q := "SELECT id, client_id, host_name, serial, os, system_data FROM agents WHERE id = $1"
	err := d.db.QueryRowContext(context.Background(), q, id).Scan(&a.ID, &a.ClientID, &a.Name, &a.Serial, &a.OS, &a.SystemData)
	if checkError(err) {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("agent '%d' does not exist: %w", id, err)
		}
		return nil, err
	}

	return a, nil
}
