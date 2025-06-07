package model

import (
    "time"
)

type Task struct {
    ID               int       `json:"id"`
    Judul            string    `json:"judul"`
    Deskripsi        string    `json:"deskripsi"`
    Prioritas        int       `json:"prioritas"`
    Status           bool      `json:"status"`
    Deadline         time.Time `json:"deadline"`
    Category         string    `json:"category"`
    TanggalPembuatan time.Time `json:"tanggalpembuatan"`
    DeleteSoft       bool      `json:"deletesoft"`
}