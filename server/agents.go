package main

import (
	"context"
	"time"
)

type agent struct {
	ID       int
	ClientID int
	Name     string
	Deleted  time.Time
}

func (svc *service) getAgents(limit, skip int) []agent {
	q := "SELECT id, client_id, hostname FROM agents WHERE id NOT IN (select id FROM deleted_agents) ORDER BY hostname asc LIMIT ? OFFSET ?"
	rows, err := svc.db.QueryContext(context.Background(), q, limit, skip)
	if checkError(err) {
		return nil
	}

	var agents []agent
	for rows.Next() {
		var a agent
		err = rows.Scan(&a.ID, &a.ClientID, &a.Name)
		if checkError(err) {
			return nil
		}

		agents = append(agents, a)
	}

	return agents
}

func (svc *service) getAgentByID(id int) []agent {
	var a agent
	q := "SELECT id, client_id, name FROM agents WHERE id = ?"
	err := svc.db.QueryRowContext(context.Background(), q, id).Scan(&a.ID, &a.ClientID, &a.Name, &a.Deleted)
	if checkError(err) {
		return nil
	}

	return a
}
