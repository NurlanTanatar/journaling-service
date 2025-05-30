package main

import (
	"context"
	"log"
	"time"
)

type Item struct {
	ID          int
	Title       string
	Completed   bool
	Criticality int
	DateStart   time.Time
	DateEnd     time.Time
}

type Incidents struct {
	Items          []Item
	Count          int
	CompletedCount int
}

func fetchIncidents() ([]Item, error) {
	var items []Item
	rows, err := DB.Query(`select id, title, completed, criticality, dateStart, dateEnd from journal order by position;`)
	if err != nil {
		return []Item{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Completed,
			&item.Criticality,
			&item.DateStart,
			&item.DateEnd)
		if err != nil {
			return []Item{}, err
		}
		items = append(items, item)
	}
	return items, nil
}

func fetchIncident(ID int) (Item, error) {
	var item Item
	err := DB.QueryRow(`select id, title, completed, criticality, dateStart, dateEnd from journal where id = ($1);`, ID).Scan(
		&item.ID,
		&item.Title,
		&item.Completed,
		&item.Criticality,
		&item.DateStart,
		&item.DateEnd)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func updateIncident(ID int, title string) (Item, error) {
	var item Item
	err := DB.QueryRow(`update journal set title = ($1) where id = ($2) returning id, title, completed, criticality, dateStart, dateEnd;`, title, ID).Scan(
		&item.ID,
		&item.Title,
		&item.Completed,
		&item.Criticality,
		&item.DateStart,
		&item.DateEnd)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func fetchCount() (int, error) {
	var count int
	err := DB.QueryRow(`select count(*) from journal;`).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func fetchCompletedCount() (int, error) {
	var count int
	err := DB.QueryRow(`select count(*) from journal where completed=true;`).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func insertIncident(title string, criticality int, dateStart time.Time, dateEnd time.Time) (Item, error) {
	count, err := fetchCount()
	if err != nil {
		return Item{}, err
	}
	var id int
	log.Printf("inserting: %v, %v, %v, %v, %v", title, count, criticality, dateStart, dateEnd)
	err = DB.QueryRow(`insert into journal (title, position, criticality, dateStart, dateEnd) values ($1, $2, $3, $4, $5) returning id`,
		title,
		count,
		criticality,
		dateStart,
		dateEnd).Scan(&id)
	if err != nil {
		return Item{}, err
	}
	item := Item{ID: id, Title: title, Completed: false, Criticality: criticality, DateStart: dateStart, DateEnd: dateEnd}
	return item, nil
}

func deleteIncident(ctx context.Context, ID int) error {
	_, err := DB.Exec("delete from journal where id = ($1)", ID)
	if err != nil {
		return err
	}
	rows, err := DB.Query("select id from journal order by position")
	if err != nil {
		return err
	}
	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for idx, id := range ids {
		_, err := DB.Exec("update journal set position = ($1) where id = ($2)", idx, id)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func orderIncidents(ctx context.Context, values []int) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, v := range values {
		_, err := tx.Exec("update journal set position = ($1) where id = ($2)", i, v)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func toggleIncident(ID int) (Item, error) {
	var item Item
	err := DB.QueryRow("update journal set completed = case when completed = true then false else true end where id = ($1) returning id, title, completed, criticality, dateStart, dateEnd", ID).Scan(
		&item.ID,
		&item.Title,
		&item.Completed,
		&item.Criticality,
		&item.DateStart,
		&item.DateEnd)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}
