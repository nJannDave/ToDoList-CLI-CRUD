package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"errors"
	"ProgramTDL/reposity"
	"ProgramTDL/model"
	"ProgramTDL/helper"
	"github.com/jackc/pgx/v5"
)

//BAB - Variabel 

const tableName = "tasks"
const obI = "id"
const obD = "deadline"
const asc = "ASC"
var c model.Task
var status string
var prioritas string
var no int
var idDS string
var idDI int
var checkId int
var exists bool
var qOpt string
var realPrioritas int //4

// BAB - Code CRUD ( Create, Read, Update, Delete )

func CreateTask() {
	db, err := reposity.GetDb()
	if err != nil {
		log.Fatalf("❌  \033[31mERROR\033[0m : %v", err)
	}
	var judul string      //1
	var deskripsi string  //2
	var qDeadline string  //5
	var kategori string   //3
	//Q - Judul
	for {
		judul = helper.InputReader("Masukkan judul Tugas yang akan ditambahkan : ")
		if judul == "" {
			fmt.Println("❌  \033[31mERROR\033[0m : Judul Tugas tidak boleh kosong! ")
			continue
		} else if len(judul) > 100 {
			fmt.Println("❌  \033[31mERROR\033[0m : Maxsimal panjang karakter judul tidak boleh lebih dari 100!")
			continue
		} else {
			break
		}
	}
	//Q - deskripsi
	for {
		deskripsi = helper.InputReader("Masukkan deskripsi Tugas yang akan ditambahkan : ")
		if deskripsi == "" {
			fmt.Println("❌  \033[31mERROR\033[0m : Deskripsi Tugas tidak boleh kosong!")
			continue
		} else {
			break
		}
	}
	//Q - Kategori
	for {
		kategori = helper.InputReader("Masukkan kategori Tugas anda! contoh = ( Pembelajaran / Hobi ) : ")
		if kategori == "" {
			fmt.Println("❌  \033[31mERROR\033[0m : Kategori Tugas tidak boleh kosong!")
			continue
		} else if len(kategori) > 20 {
			fmt.Println("❌  \033[31mERROR\033[0m : Panjang karakter kategori Tugas, maxsimal tidak boleh lebih dari 20!")
			continue
		} else {
			break
		}
	}
	//Q - Prioritas
	for {
		qPrioritas := helper.InputReader("Masukkan prioritas Tugas anda! ( Penting / Lumayan Penting / Tidak Penting ) : ")
		if qPrioritas == "" {
			fmt.Println("❌  \033[31mERROR\033[0m : Prioritas tidak boleh kosong!")
			continue
		} else if strings.ToLower(qPrioritas) == "penting" {
			realPrioritas = 1
			break
		} else if strings.ToLower(qPrioritas) == "lumayan penting" {
			realPrioritas = 2
			break
		} else if strings.ToLower(qPrioritas) == "tidak penting" {
			realPrioritas = 3
			break
		} else {
			fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi prioritas Tugas yg tersedia! ( Penting / Lumayan Penting / Tidak Penting )")
			continue
		}
	}
	//Q - Deadline
	for {
		qDl := helper.InputReader("Masukkan deadline Tugas! contoh = ( YYYY-MM-DD, YYYY/MM/DD atau 2025-04-02, 2025/04/02 ) : ")
		if qDl == "" {
			fmt.Println("❌  \033[31mERROR\033[0m : Deadline tidak boleh kosong!")
			continue
		}
		qDeadline = strings.Replace(qDl, "/", "-", -1)
		realDeadline, err := time.Parse("2006-01-02", qDeadline)
		if err != nil {
			fmt.Println("\033[31mERROR\033[0m : Format tanggal salah! Gunakan format seperti ini YYYY-MM-DD atau YYYY/MM/DD")
			continue
		}
		if realDeadline.Before(time.Now().Truncate(24 * time.Hour)) {
			fmt.Println("❌  \033[31mERROR\033[0m : Deadline tidak boleh masa lalu!")
			continue
		}
		//Bab - Query
		query := fmt.Sprintf("INSERT INTO %s (judul, deskripsi, category, prioritas,deadline) VALUES ($1, $2, $3, $4, $5)", tableName)
		_, err = db.Exec(context.Background(), query, judul, deskripsi, kategori, realPrioritas, realDeadline)
		if err != nil {
			log.Printf("\033[31mERROR\033[0m : %v", err)
			fmt.Printf("❌  \033[31mERROR\033[0m : Gagal menambahkan Tugas! : %v", err)
			helper.EmptySpace()
			retry := helper.InputReader("Coba lagi? (y/n): ")
			if strings.ToLower(retry) == "y" {
				continue
			} else {
				fmt.Println("🚫 Operasi dibatalkan")
				return
			}
		} else {
			helper.EmptySpace()
			fmt.Println("✅ Tugas \033[32mBerhasil\033[0m ditambahkan!")
			helper.EmptySpace()
			helper.DecorTwo()
			 
		}
		break
	}
}

func ReadTask() {
	db, err := reposity.GetDb()
	if err != nil {
		log.Fatalf("❌  \033[31mERROR\033[0m : %v", err)
	}
	helper.EmptySpace()
	helper.DecorTwo()
	helper.EmptySpace()
	fmt.Println("Opsi yang tersedia : ")
	helper.EmptySpace()
	helper.DecorTwo()
	helper.EmptySpace()
	helper.SpecialDecorThree()
	helper.EmptySpace()
	var opRead string
	for {
		helper.DecorTwo()
		helper.EmptySpace()
		opRead = helper.InputReader("Masukkan opsi yang diinginkan! : ")
		if opRead == "" {
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("❌  \033[31mERROR\033[0m : Opsi tidak boleh kosong, harus memilih!")
			helper.EmptySpace()
			continue
		} else if opRead == "1" {
		     
			query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id, judul, category, prioritas, status FROM %s WHERE deletesoft = $1", obI, asc, tableName)
			row, err := db.Query(context.Background(), query, false)
			if err != nil {
				log.Printf("\033[31mERROR\033[0m : %v", err)
				fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
				helper.EmptySpace()
				retry := helper.InputReader("Coba lagi? (y/n): ")
				if strings.ToLower(retry) == "y" {
					continue
				} else {
					fmt.Println("🚫 Operasi dibatalkan")
					return
				}
			}
			defer row.Close()
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("Daftar Tugas : ")
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			for row.Next() {
				if err := row.Scan(&no ,&c.ID , &c.Judul, &c.Category, &c.Prioritas, &c.Status); err != nil {
					log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
					continue
				}
				if c.Status == false {
					status = "\033[31mBelum Selesai\033[0m"
				} else {
					status = "\033[32mSelesai\033[0m"
				}
				if c.Prioritas == 1 {
					prioritas = "\033[31mPenting\033[0m"
				} else if c.Prioritas == 2 {
					prioritas = "\033[33mLumayan Penting\033[0m"
				} else if c.Prioritas == 3 {
					prioritas = "\033[32mTidak Penting\033[0m"
				}
				fmt.Printf("No %v. [ID : %d] - %s - %s - %s - %s\n",no , c.ID, c.Judul, c.Category, prioritas, status)
			}
			helper.EmptySpace()
			helper.DecorTwo()
			break
		} else if opRead == "2" {
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("Opsi yang tersedia berdasarkan Filter : ")
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			helper.SpecialDecorFourht()
			helper.EmptySpace()
			var opReadF string
			for {
				helper.DecorTwo()
				helper.EmptySpace()
				opReadF = helper.InputReader("Masukkan opsi yang ingin digunakan! ( 1 / 2 / 3 / 4 ) : ")
				helper.EmptySpace()
				helper.DecorTwo()
				if opReadF == "" {
					helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Opsi tidak boleh kosong, harus memilih!")
					helper.EmptySpace()
					continue
				} else if opReadF == "1" {
				     
					query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s $s ) AS no, id, judul, deskripsi, status FROM %s WHERE status = $1 AND deletesoft = $2", obI, asc, tableName)
					row, err := db.Query(context.Background(), query, true, false)
					if err != nil {
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					}
					defer row.Close()
					helper.EmptySpace()
					fmt.Println("Daftar Tugas berdasarkan Status = Selesai : ")
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					for row.Next() {
						if err := row.Scan(&no ,&c.ID, &c.Judul, &c.Deskripsi, &c.Status); err != nil {
							log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
							continue
						}
						if c.Status == false {
							status = "\033[31mBelum Selesai\033[0m"
						} else {
							status = "\033[31mSelesai\033[0m"
						}
						fmt.Printf("No %v. [%d] - %s - %s - %s\n", no, c.ID, c.Judul, c.Deskripsi, status)
					}
					helper.EmptySpace()
					helper.DecorTwo()
					break
				} else if opReadF == "2" {
				     
					query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id, judul, deskripsi, status FROM %s WHERE status = $1 AND deletesoft = $2",obI, asc, tableName)
					row, err := db.Query(context.Background(), query, false, false)
					if err != nil {
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					}
					defer row.Close()
					helper.EmptySpace()
					fmt.Println("Daftar Tugas berdasarkan Status = Belum Selesai : ")
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					for row.Next() {
						if err := row.Scan(&no, &c.ID, &c.Judul, &c.Deskripsi, &c.Status); err != nil {
							log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
							continue
						}
						if c.Status == false {
							status = "\033[31mBelum Selesai\033[0m"
						} else {
							status = "\033[31mSelesai\033[0m"
						}
						fmt.Printf("No %v. [%d] - %s - %s - %s\n", no, c.ID, c.Judul, c.Deskripsi, status)
					}
					helper.EmptySpace()
					helper.DecorTwo()
					break
				} else if opReadF == "3" {
				     
					query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id, judul, deskripsi, prioritas FROM %s WHERE prioritas = $1 AND deletesoft = $2", obI, asc, tableName)
					row, err := db.Query(context.Background(), query, 1, false)
					if err != nil {
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil Tugas : %v", err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					}
					defer row.Close()
					helper.EmptySpace()
					fmt.Println("Daftar Tugas berdasarkan Prioritas = Penting")
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					for row.Next() {
						if err := row.Scan(&no, &c.ID, &c.Judul, &c.Deskripsi, &c.Prioritas); err != nil {
							log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
							continue
						}
						if c.Prioritas == 1 {
							prioritas = "\033[31mPenting\033[0m"
						} else if c.Prioritas == 2 {
							prioritas = "\033[33mLumayan Penting\033[0m"
						} else if c.Prioritas == 3 {
							prioritas = "\033[32mTidak Penting\033[0m"
						}
						fmt.Printf("No %v. [%d] - %s - %s - %s\n", no, c.ID, c.Judul, c.Deskripsi, prioritas)
					}
					helper.EmptySpace()
					helper.DecorTwo()
					break
				} else if opReadF == "4" {
				     
					query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s ) AS no, id, judul, deskripsi, deadline FROM %s WHERE status = $1 AND deadline >= CURRENT_DATE AND deletesoft = $2", obD, asc, tableName)
					row, err := db.Query(context.Background(), query, false, false)
					if err != nil {
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil Tugas : %v", err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					}
					helper.EmptySpace()
					fmt.Println("Daftar Tugas berdasarkan Deadline Terdekat : ")
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					for row.Next() {
						if err := row.Scan(&no, &c.ID, &c.Judul, &c.Deskripsi, &c.Deadline); err != nil {
							log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
							continue
						}
						sisaHari := int(time.Until(c.Deadline).Hours() / 24)
						fmt.Printf("No %v. %s (ID : %d) - %s - \033[33m%s\033[0m (\033[31m%d hari lagi\033[0m)\n", no, c.Judul, c.ID, c.Deskripsi, c.Deadline.Format("02 Jan 2006"), sisaHari)
					}
					helper.EmptySpace()
					helper.DecorTwo()
					break
				} else {
					helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia!")
					helper.EmptySpace()
					continue
				}
			}
			break
		} else {
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia!")
			helper.EmptySpace()
			continue
		}
	}
}

func DeleteTask() {
	db, err := reposity.GetDb()
	if err != nil {
		log.Fatalf("❌  \033[31mERROR\033[0m : %v", err)
	}
	var qConfirm string
	var deletesoft bool
	helper.EmptySpace()
	helper.DecorTwo()
	helper.EmptySpace()
	fmt.Println("Opsi yang tersedia : ")
	helper.EmptySpace()
	helper.DecorTwo()
	helper.EmptySpace()
	helper.SpecialDecorFive()
	helper.EmptySpace()
	for {
		helper.DecorTwo()
		helper.EmptySpace()
		qOpt = helper.InputReader("Masukkan opsi Delete yang akan digunakan! ( 1 / 2 ) : ")
		if qOpt == "" {
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("❌  \033[31mERROR\033[0m : Opsi Delete tidak boleh kosong, harus memilih!")
			helper.EmptySpace()
			continue
		} else if qOpt == "1" {
			helper.ReadTaskH()
			helper.EmptySpace()
			helper.DecorTwo()
			for {
				helper.EmptySpace()
				idDS = helper.InputReader("Masukkan NO Tugas yang akan di SoftDelete ( Tidak Permanen ) : ")
				if idDS == "" {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Printf("❌  \033[31mERROR\033[0m : Tidak boleh kosong!")
					helper.EmptySpace()
					helper.EmptySpace()
					helper.DecorTwo()
					continue
				}
				idDI, err = strconv.Atoi(idDS)
				if err != nil {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Printf("❌  \033[31mERROR\033[0m : Harus berupa angka! : %v", err)
					helper.EmptySpace()
					helper.EmptySpace()
					helper.DecorTwo()
					continue
				}
				 
				query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
				break
			}
			for {
				helper.EmptySpace()
				helper.DecorTwo()
				helper.EmptySpace()
				qConfirm = helper.InputReader("Apakah anda yakin ingin menghapus ( SoftDelete / Tidak Permanen ) tugas tersebut? ( y / n ) : ")
				if qConfirm == "" {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Pilih salah satu = ( y / n ), tidak boleh kosong!")
					continue
				} else if strings.ToLower(qConfirm) == "n" {
					deletesoft = false
					break
				} else if strings.ToLower(qConfirm) == "y" {
					deletesoft = true
					query := fmt.Sprintf("UPDATE %s SET deletesoft = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id from %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
					_, err := db.Exec(context.Background(), query, deletesoft, false, idDI)
					if err != nil {
					    helper.EmptySpace()
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-SoftDelete Tugas dengan NO %v : %v", idDI, err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
						    helper.EmptySpace()
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					} else {
						helper.EmptySpace()
						helper.DecorTwo()
						helper.EmptySpace()
						fmt.Printf("✅ Tugas dengan NO %v, \033[32mBerhasil\033[0m dihapus ( SoftDelete / Tidak Permanen )!", idDI)
						helper.EmptySpace()
						helper.EmptySpace()
						helper.DecorTwo()
						helper.EmptySpace()
					}
				} else {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia, ( y / n )!")
					continue
				}
				break
			}
			break
		} else if qOpt == "2" {
		     
			query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id, judul, category, prioritas, status FROM %s", obI, asc, tableName)
			row, err := db.Query(context.Background(), query)
			if err != nil {
			    helper.EmptySpace()
				log.Printf("\033[31mERROR\033[0m : %v", err)
				fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
				helper.EmptySpace()
				retry := helper.InputReader("Coba lagi? (y/n): ")
				if strings.ToLower(retry) == "y" {
					continue
				} else {
					helper.EmptySpace()
					fmt.Println("🚫 Operasi dibatalkan")
					return
				}
			}
			defer row.Close()
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("Daftar Tugas : ")
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			for row.Next() {
				if err := row.Scan(&no, &c.ID, &c.Judul, &c.Category, &c.Prioritas, &c.Status); err != nil {
					log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
					continue
				}
				if c.Status == false {
					status = "Belum Selesai"
				} else {
					status = "Selesai"
				}
				if c.Prioritas == 1 {
					prioritas = "\033[31mPenting\033[0m"
				} else if c.Prioritas == 2 {
					prioritas = "\033[33mLumayan Penting\033[0m"
				} else if c.Prioritas == 3 {
					prioritas = "\033[32mTidak Penting\033[0m"
				}
				fmt.Printf("No %v. [%d] - %s - %s - %s - %s\n", no, c.ID, c.Judul, c.Category, prioritas, status)
			}
			helper.EmptySpace()
			helper.DecorTwo()
			for {
				helper.EmptySpace()
				idDS = helper.InputReader("Masukkan NO Tugas yang akan di HardDelete ( Permanen ) : ")
				if idDS == "" {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Printf("❌  \033[31mERROR\033[0m : Tidak boleh kosong!")
					helper.EmptySpace()
					helper.EmptySpace()
					helper.DecorTwo()
					continue
				}
				idDI, err = strconv.Atoi(idDS)
				if err != nil {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Printf("❌  \033[31mERROR\033[0m : Harus berupa angka! : %v", err)
					helper.EmptySpace()
					helper.EmptySpace()
					helper.DecorTwo()
					continue
				}
				query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s) AS subquery WHERE no = $1 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
				break
			}
			for {
				helper.EmptySpace()
				helper.DecorTwo()
				helper.EmptySpace()
				qConfirm = helper.InputReader("Apakah anda yakin ingin menghapus ( HardDelete / Permanen ) tugas tersebut? ( y / n ) : ")
				if qConfirm == "" {
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Pilih salah satu = ( y / n ), tidak boleh kosong!")
					continue
				} else if strings.ToLower(qConfirm) == "n" {
					break
				} else if strings.ToLower(qConfirm) == "y" {
					query := fmt.Sprintf("DELETE FROM %s WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s) AS subquery WHERE no = $1);", tableName, obI, asc, tableName)
					_, err := db.Exec(context.Background(), query, idDI)
					if err != nil {
					    helper.EmptySpace()
						log.Printf("\033[31mERROR\033[0m : %v", err)
						fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-HardDelete Tugas dengan NO %v : %v", idDI, err)
						helper.EmptySpace()
						retry := helper.InputReader("Coba lagi? (y/n): ")
						if strings.ToLower(retry) == "y" {
							continue
						} else {
						    helper.EmptySpace()
							fmt.Println("🚫 Operasi dibatalkan")
							return
						}
					} else {
						helper.EmptySpace()
						helper.DecorTwo()
						helper.EmptySpace()
						fmt.Printf("✅ Tugas dengan NO %v, \033[32mBerhasil\033[0m dihapus ( HardDelete / Permanen )!", idDI)
						helper.EmptySpace()
						helper.EmptySpace()
						helper.DecorTwo()
						helper.EmptySpace()
					}
					break
				} else {
				    helper.EmptySpace()
				    helper.DecorTwo()
				    helper.EmptySpace()
					fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia! ( y / n )")
					continue
				}
			}
		} else {
			helper.EmptySpace()
			helper.DecorTwo()
			helper.EmptySpace()
			fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi Delete yang tersedia!")
			helper.EmptySpace()
			continue
		}
		break
	}
}

func UpdateTask() {
    db, err := reposity.GetDb()
    if err != nil {
        log.Fatalf("❌  \033[31mERROR\033[0m : %v", err)
    }
    helper.EmptySpace()
    helper.DecorTwo()
	helper.EmptySpace()
	fmt.Println("Opsi yang tersedia : ")
	helper.EmptySpace()
	helper.DecorTwo()
	helper.EmptySpace()
	helper.SpecialDecorSix()
	helper.EmptySpace()
	for {
	    helper.DecorTwo()
	    helper.EmptySpace()
	    qOptUp := helper.InputReader("Masukkan opsi Update yang akan digunakan! ( 1 / 2 / 3 / 4 / 5 / 6 / 7 ) : ")
	    if qOptUp == "" {
	        helper.EmptySpace()
	        helper.DecorTwo()
	        helper.EmptySpace()
	        fmt.Printf("❌  \033[31mERROR\033[0m : Opsi update harus memilih, tidak boleh kosong!")
	        helper.EmptySpace()
	        helper.EmptySpace()
	        continue
	    } else if qOptUp == "1" {
	        helper.ReadTaskH()
			helper.EmptySpace()
			helper.DecorTwo()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Masukkan no Tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					fmt.Printf("❌  \033[31mERROR\033[0m : Harus berupa angka! : %v", err)
					helper.EmptySpace()
					helper.EmptySpace()
					helper.DecorTwo()
					continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			    helper.EmptySpace()
			    helper.DecorTwo()
			    helper.EmptySpace()
			}
			for {
    			judulBaru := helper.InputReader("Masukkan Judul baru untuk menggantikan Judul lama! : ")
    			if judulBaru == "" {
    			    helper.EmptySpace()
    			    helper.DecorTwo()
    			    helper.EmptySpace()
    			    fmt.Printf("❌  \033[31mERROR\033[0m : Masukkan judul baru, tidak boleh kosong!")
    			    continue
    			} else if len(judulBaru) > 100 {
    			    helper.EmptySpace()
    			    helper.DecorTwo()
    			    helper.EmptySpace()
    			    fmt.Println("❌  \033[31mERROR\033[0m : Maxsimal panjang karakter judul tidak boleh lebih dari 100!")
    			    continue
    			} else {
    			    query := fmt.Sprintf("UPDATE %s SET judul = $1 WHERE id = (SELECT id FROM  (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
    			    _, err := db.Exec(context.Background(), query, judulBaru, false, idDI)
    			    if err != nil {
    			        log.Printf("\033[31mERROR\033[0m : %v", err)
    			        fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update judul Tugas dengan No %v! : %v", idDI, err)
    			        helper.EmptySpace()
    				    retry := helper.InputReader("Coba lagi? (y/n): ")
    				    if strings.ToLower(retry) == "y" {
    						continue
    					} else {
    						helper.EmptySpace()
    					    fmt.Println("🚫 Operasi dibatalkan")
    						return
    					}
    			    } else {
    			        helper.EmptySpace()
    			        helper.DecorTwo()
    			        helper.EmptySpace()
    			        fmt.Printf("✅ Berhasil meng-Update judul Tugas dengan No %v!", idDI)
    			        helper.EmptySpace()
    			        helper.EmptySpace()
    			        helper.DecorTwo()
    			        break
    			    }
    			}
			}
			break
	    } else if qOptUp == "2" {
	        helper.ReadTaskH()
			helper.EmptySpace()
			helper.DecorTwo()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Masukkan no tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			for {
			    helper.EmptySpace()
			    helper.DecorTwo()
			    helper.EmptySpace()
			    deskripsiBaru := helper.InputReader("Masukkan Deskripsi baru untuk menggantikan Deskripsi lama! : ")
			    if deskripsiBaru == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Deskripsi tidak boleh kosong, harus mengisi!")
			        helper.EmptySpace()
			        continue
			    } else {
			        query := fmt.Sprintf("UPDATE %s SET deskripsi = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
			        _, err := db.Exec(context.Background(), query, deskripsiBaru, false, idDI)
			        if err != nil {
			            log.Printf("\033[31mERROR\033[0m : %v", err)
			            fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update judul Tugas dengan No %v! : %v", idDI, err)
    			        helper.EmptySpace()
    			        helper.EmptySpace()
    				    retry := helper.InputReader("Coba lagi? (y/n): ")
    				    if strings.ToLower(retry) == "y" {
    						continue
    					} else {
    						helper.EmptySpace()
    					    fmt.Println("🚫 Operasi dibatalkan")
    						return
    					}
			        } else {
			            helper.EmptySpace()
    			        helper.DecorTwo()
    			        helper.EmptySpace()
    			        fmt.Printf("✅ Berhasil meng-Update judul Tugas dengan No %v!", idDI)
    			        helper.EmptySpace()
    			        helper.EmptySpace()
    			        helper.DecorTwo()
    			        break
			        }
			    }
			}
			break
	    } else if qOptUp == "3" {
	        helper.ReadTaskH()
			helper.EmptySpace()
			helper.DecorTwo()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Masukkan no tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			for {
			    helper.EmptySpace()
			    helper.DecorTwo()
			    helper.EmptySpace()
			    kategoriBaru := helper.InputReader("Masukkan Kategori baru untuk menggantikan Kategori lama! : ")
			    if kategoriBaru == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Kategori tidak boleh kosong, harus mengisi!")
			        helper.EmptySpace()
			        continue
			    } else if len(kategoriBaru) > 20 {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Maxsimal panjang karakter Kategori adalah 20, tidak boleh lebih!")
			        helper.EmptySpace()
			        continue
			    } else {
			        query := fmt.Sprintf("UPDATE %s SET category = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
			        _, err := db.Exec(context.Background(), query, kategoriBaru, false, idDI)
			        if err != nil {
			            log.Printf("\033[31mERROR\033[0m : %v", err)
			            fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update Kategori Tugas dengan No %v! : %v", idDI, err)
    			        helper.EmptySpace()
    			        helper.EmptySpace()
    				    retry := helper.InputReader("Coba lagi? (y/n): ")
    				    if strings.ToLower(retry) == "y" {
    						continue
    					} else {
    						helper.EmptySpace()
    					    fmt.Println("🚫 Operasi dibatalkan")
    						return
    					}
			        } else {
			            helper.EmptySpace()
    			        helper.DecorTwo()
    			        helper.EmptySpace()
    			        fmt.Printf("✅ Berhasil meng-Update Kategori Tugas dengan No %v!", idDI)
    			        helper.EmptySpace()
    			        helper.EmptySpace()
    			        helper.DecorTwo()
    			        break
			        }
			    }
			}
			break
	    } else if qOptUp == "4" {
	        helper.ReadTaskH()
			helper.EmptySpace()
			helper.DecorTwo()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Masukkan no tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			for {
			    helper.EmptySpace()
			    helper.DecorTwo()
			    helper.EmptySpace()
			    deadlineBaru := helper.InputReader("Masukkan Deadline baru untuk menggantikan Deadline lama! contoh = ( YYYY/MM/DD atau YYYY-MM-DD ) : ")
			    if deadlineBaru == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Deadline tidak boleh kosong, harus mengisi!")
			        helper.EmptySpace()
			        continue
			    }
			    deadlineBaru = strings.Replace(deadlineBaru, "/", "-", -1)
		        realDeadline, err := time.Parse("2006-01-02", deadlineBaru)
		        if err != nil {
		            helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Println("❌  \033[31mERROR\033[0m : Format tanggal salah! Gunakan format seperti ini YYYY-MM-DD atau YYYY/MM/DD")
			        continue
		        }
		        if realDeadline.Before(time.Now().Truncate(24 * time.Hour)) {
		            helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Println("❌  \033[31mERROR\033[0m : Deadline tidak boleh masa lalu!")
			        continue
		        }
			    query := fmt.Sprintf("UPDATE %s SET deadline = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obD, asc, tableName)
			    _, err = db.Exec(context.Background(), query, realDeadline, false, idDI)
			    if err != nil {
			        log.Printf("\033[31mERROR\033[0m : %v", err)
			        fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update Deadline Tugas dengan No %v! : %v", idDI, err)
    			    helper.EmptySpace()
    			    helper.EmptySpace()
    			    retry := helper.InputReader("Coba lagi? (y/n): ")
    			    if strings.ToLower(retry) == "y" {
    			        continue
    				} else {
    					helper.EmptySpace()
    					fmt.Println("🚫 Operasi dibatalkan")
    					return
    				}
			    } else {
			        helper.EmptySpace()
    			    helper.DecorTwo()
    			    helper.EmptySpace()
    			    fmt.Printf("✅ Berhasil meng-Update Deadline Tugas dengan No %v!", idDI)
    			    helper.EmptySpace()
    			    helper.EmptySpace()
    			    helper.DecorTwo()
    			    break
			    }
		    }
		break
	    } else if qOptUp == "5" {
	        helper.ReadTaskH()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Masukkan no tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			for {
			    helper.EmptySpace()
			    helper.DecorTwo()
			    helper.EmptySpace()
			    prioritasBaru := helper.InputReader("Masukkan tingkat Prioritas baru! ( Penting / Lumayan Penting / Tidak Penting ) : ")
			    if prioritasBaru == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Tingikat Priotitas tidak boleh kosong, harus mengisi!")
			        helper.EmptySpace()
			        continue
			    } else if strings.ToLower(prioritasBaru) == "penting" {
			        realPrioritas = 1
			        break
			    } else if strings.ToLower(prioritasBaru) == "lumayan penting" { 
			       realPrioritas = 2
			       break
			    } else if strings.ToLower(prioritasBaru) == "tidak penting" {
			        realPrioritas = 3
			        break
			    } else {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Pilih kategori yang tersedia! ( Penting / Lumayan Penting / Tidak Penting )")
			        helper.EmptySpace()
			        continue
			    }
			}
			query := fmt.Sprintf("UPDATE %s SET prioritas = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
			_, err = db.Exec(context.Background(), query, realPrioritas, false, idDI)
			    if err != nil {
			        log.Printf("\033[31mERROR\033[0m : %v", err)
			        fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update Prioritas Tugas dengan No %v! : %v", idDI, err)
    			    helper.EmptySpace()
    			    helper.EmptySpace()
    			    retry := helper.InputReader("Coba lagi? (y/n): ")
    			    if strings.ToLower(retry) == "y" {
    					continue
    				} else {
    					helper.EmptySpace()
    					fmt.Println("🚫 Operasi dibatalkan")
    					return
    				}
			    } else {
			        helper.EmptySpace()
    			    helper.DecorTwo()
    			    helper.EmptySpace()
    			    fmt.Printf("✅ Berhasil meng-Update Prioritas Tugas dengan No %v!", idDI)
    			    helper.EmptySpace()
    			    helper.EmptySpace()
    			    helper.DecorTwo()
    			    break
                }
		    break
	    } else if qOptUp == "6" {
	        helper.ReadTaskH()
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Update! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Masukkan no tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("\033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $1) AS subquery WHERE no = $2 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, false, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			for {
			    helper.EmptySpace()
			    qStatus := helper.InputReader("Tandai Status Tugas sudah selesai!? ( y / n )")
			    if qStatus == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Tandai Tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    } else if strings.ToLower(qStatus) == "y" {
			        stts := true
			        query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s WHERE deletesoft = $2) AS subquery WHERE no = $3);", tableName, obI, asc, tableName)
			        _, err := db.Exec(context.Background(), query, stts, false, idDI)
			        if err != nil {
        			    log.Printf("\033[31mERROR\033[0m : %v", err)
        			    fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update Prioritas Tugas dengan No %v! : %v", idDI, err)
            		    helper.EmptySpace()
            		    helper.EmptySpace()
            			retry := helper.InputReader("Coba lagi? (y/n): ")
            			if strings.ToLower(retry) == "y" {
            				continue
            			} else {
            				helper.EmptySpace()
            				fmt.Println("🚫 Operasi dibatalkan")
            				return
            			}
        		    } else {
        			    helper.EmptySpace()
            			helper.DecorTwo()
            			helper.EmptySpace()
            			fmt.Printf("✅ Berhasil meng-Update Prioritas Tugas dengan No %v!", idDI)
            			helper.EmptySpace()
                        helper.EmptySpace()
            			helper.DecorTwo()
                    }
                    break
			    } else if strings.ToLower(qStatus) == "n" {
			        break
			    } else {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Pilih salah satu opsi yaitu ( y / n )!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			}
			break
	    } else if qOptUp == "7" {
	        rDB := true
	        rDb2 := false
	        query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s ) AS no, id, judul, category, prioritas, status FROM %s WHERE deletesoft = $1", obI, asc, tableName)
            row, err := db.Query(context.Background(), query, rDB)
            if err != nil {
            	helper.EmptySpace()
            	log.Printf("\033[31mERROR\033[0m : %v", err)
            	fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
            	helper.EmptySpace()
            	retry := helper.InputReader("Coba lagi? (y/n): ")
            	if strings.ToLower(retry) == "y" {
            		continue
            	} else {
            		helper.EmptySpace()
            		fmt.Println("🚫 Operasi dibatalkan")
            		return
            	}
            }
            defer row.Close()
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("Daftar Tugas : ")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            for row.Next() {
            	if err := row.Scan(&no, &c.ID, &c.Judul, &c.Category, &c.Prioritas, &c.Status); err != nil {
            		log.Printf("❌  \033[31mERROR\033[0m : Gagal membaca data : %v", err)
            		continue
            	}
            	if c.Status == false {
            		status = "Belum Selesai"
            	} else {
            		status = "Selesai"
            	}
            	if c.Prioritas == 1 {
            		prioritas = "\033[31mPenting\033[0m"
            	} else if c.Prioritas == 2 {
            		prioritas = "\033[33mLumayan Penting\033[0m"
            	} else if c.Prioritas == 3 {
            		prioritas = "\033[32mTidak Penting\033[0m"
            	}
            	fmt.Printf("No %v. [%d] - %s - %s - %s - %s\n", no, c.ID, c.Judul, c.Category, prioritas, status)
            }
			for {
			    helper.EmptySpace()
			    idDS = helper.InputReader("Masukkan No Tugas yang akan di Restore / Dikembalikan! : ")
			    if idDS == "" {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌  \033[31mERROR\033[0m : Masukkan No tugas, tidak boleh kosong!")
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    idDI, err = strconv.Atoi(idDS)
			    if err != nil {
			        helper.EmptySpace()
			        helper.DecorTwo()
			        helper.EmptySpace()
			        fmt.Printf("❌ \033[31mERROR\033[0m : Harus berupa angka! : %v", err)
			        helper.EmptySpace()
			        helper.EmptySpace()
			        helper.DecorTwo()
			        continue
			    }
			    query := fmt.Sprintf("SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s) AS subquery WHERE no = $1 ", obI, asc, tableName)
				err = db.QueryRow(context.Background(), query, idDI).Scan(&checkId)
				if err == nil {
				    exists = true
				    break
				} else if errors.Is(err, pgx.ErrNoRows) {
				    exists = false
				}
				if !exists {
					fmt.Printf("❌  \033[31mERROR\033[0m : Tugas dengan NO %v tidak ditemukan!\n", idDI)
					helper.EmptySpace()
					helper.DecorTwo()
					helper.EmptySpace()
					retry := helper.InputReader("Coba lagi? (y/n): ")
					if strings.ToLower(retry) == "y" {
						continue
					} else {
					    helper.EmptySpace()
						fmt.Println("🚫 Operasi dibatalkan")
						return
					}
				}
			}
			query = fmt.Sprintf("UPDATE %s SET deletesoft = $1 WHERE id = (SELECT id FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s %s) AS no, id FROM %s) AS subquery WHERE no = $2);", tableName, obI, asc, tableName)
			_, err = db.Exec(context.Background(), query, rDb2, idDI)
			if err != nil {
			    log.Printf("\033[31mERROR\033[0m : %v", err)
			    fmt.Printf("❌  \033[31mERROR\033[0m : Gagal meng-Update Prioritas Tugas dengan No %v! : %v", idDI, err)
    		    helper.EmptySpace()
    		    helper.EmptySpace()
    			retry := helper.InputReader("Coba lagi? (y/n): ")
    			if strings.ToLower(retry) == "y" {
    				continue
    			} else {
    				helper.EmptySpace()
    				fmt.Println("🚫 Operasi dibatalkan")
    				return
    			}
		    } else {
			    helper.EmptySpace()
    			helper.DecorTwo()
    			helper.EmptySpace()
    			fmt.Printf("✅ Berhasil meng-Update Prioritas Tugas dengan No %v!", idDI)
    			helper.EmptySpace()
                helper.EmptySpace()
    			helper.DecorTwo()
            }
            break
	    } else {
	        helper.EmptySpace()
	        helper.DecorTwo()
	        helper.EmptySpace()
	        fmt.Printf("❌  \033[31mERROR\033[0m : Pilih opsi update yang tersedia ( 1 / 2 / 3 / 4 / 5 / 6 / 7 )!")
	        helper.EmptySpace()
	        helper.EmptySpace()
	        continue
	    }
    }
}