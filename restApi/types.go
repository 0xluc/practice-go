package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectId    int64     `json:"projectId"`
	AssignedToId int64     `json:"assignedTo"`
	CreatedAt    time.Time `json:"createdAt"`
}

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}
