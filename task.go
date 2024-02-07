package main

import "context"

type Item struct {
	ID        int
	Title     string
	Completed bool
}

type Tasks struct {
	Items          []int
	Count          int
	CompletedCount int
}

func fetchTasks() ([]Item, error) {
	var items []Item
	rows, err := DB.Query(`select id, title, completed from tasks order by position;`)
	if err != nil {
		return []Item{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Title, &item.Completed)
		if err != nil {
			return []Item{}, err
		}
		items = append(items, item)
	}
	return items, nil
}

func fetchTask(ID int) (Item, error) {
	var item Item
	err := DB.QueryRow(`select id, title, completed from tasks where id = (?);`, ID).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func updateTask(ID int, title string) (Item, error) {
	var item Item
	err := DB.QueryRow(`update tasks set title = (?) where id = (?) returning id, title, completed;`, title, ID).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return Item{}, err
	}
	return item, nil
}

func fetchCount() (int, error) {
	var count int
	err := DB.QueryRow(`select count(*) from tasks;`).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func insertTask(title string) (Item, error) {
	count, err := fetchCount()
	if err != nil {
		return Item{}, err
	}
	var id int
	err = DB.QueryRow(`insert into tasks (title, position) values (?, ?) returning id;`, title, count).Scan(&id)
	if err != nil {
		return Item{}, err
	}
	item := Item{ID: id, Title: title, Completed: false}
	return item, nil
}

func deleteTask(ctx context.Context, ID int) error {
	_, err := DB.Exec(`delete from tasks where id = (?);`, ID)
	if err != nil {
		return err
	}
	return nil
}
