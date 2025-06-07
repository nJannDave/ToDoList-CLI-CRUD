package main
import (
    "ProgramTDL/reposity"
    "ProgramTDL/controllers"
    "ProgramTDL/helper"
    "fmt"
    "log"
    "strings"
    "os"
)

//Bab - Fungsi

func mainMenuTdl (userName string) {
    helper.EmptySpace()
    fmt.Println("Hello, ", userName, "!!")
    fmt.Println("Selamat datang di menu utama...")
    helper.EmptySpace()
    helper.SpecialDecorOne()
    helper.EmptySpace()
    fmt.Println("Opsi / Pilihan yang tersedia!")
    helper.EmptySpace()
    helper.SpecialDecorTwo()
    helper.EmptySpace()
    helper.DecorTwo()
    helper.EmptySpace()
    var option string
    for {
        option = helper.InputReader("Pilih opsi yang akan digunakan! ( 1 / 2 / 3 / 4 / 5 ) : ")
        if option == "" {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi, tidak boleh kosong!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        } else if option == "5" {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("Baiklah...")
            fmt.Println("Terimakasih telah menggunakan, ", userName, "!!")
            fmt.Println("🔴  Program akan berhenti!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            break
        } else if option == "1" {
            controllers.CreateTask()
            break
        } else if option == "2" {
            controllers.ReadTask()
            break
        } else if option == "3" {
            controllers.UpdateTask()
            break
        } else if option == "4" {
            controllers.DeleteTask()
            break
        } else {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        }
    }
}

func sayHelloAndValidation(userName string) {
    helper.EmptySpace()
    fmt.Println("Hello, ", userName, "!!")
    helper.EmptySpace()
    helper.DecorTwo()
    helper.EmptySpace()
    fmt.Println("Selamat datang di program...")
    helper.EmptySpace()
    helper.SpecialDecorOne()
    helper.EmptySpace()
    helper.DecorTwo()
    helper.EmptySpace()
    var validation string
    for {
        validation = helper.InputReader("Jadi... apakah anda ingin menggunakan Program ini? ( y / n ) : ")
        if validation == "" {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Pilihan tidak boleh kosong!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        } else if strings.ToLower(validation) == "n" {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("Baiklah...")
            fmt.Println("Terimakasih telah menggunakan, ", userName, "!!")
            fmt.Println("🔴  Program akan berhenti!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            break
        } else if strings.ToLower(validation) == "y" {
            mainMenuTdl(userName)
            break
        } else {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Pilih opsi yang tersedia! Yaitu = y / n ")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        }
    }
}

func getUserName() {
    helper.DecorOne()
    helper.EmptySpace()
    var userName string
    for {
        userName = helper.InputReader("Masukkan Nama / Username anda terlebih dahulu! ( Wajib diisi! ) : ")
        if userName == "" {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Nama wajib diisi, tidak boleh kosong!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        } else if len(userName) > 20 {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("❌  \033[31mERROR\033[0m : Panjang karakter nama minimal tidak lebih dari 20 karakter!")
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            continue
        } else {
            helper.EmptySpace()
            helper.DecorTwo()
            helper.EmptySpace()
            fmt.Println("✅ Terimakasih sudah mengisi!")
            helper.EmptySpace()
            helper.DecorTwo()
            break
        }
    }
    sayHelloAndValidation(userName)
}

//Bab - Fungsi Utama
func main() {
    if err := reposity.ConnectDb(); err != nil {
        helper.EmptySpace()
        helper.DecorTwo()
        helper.EmptySpace()
        log.Fatalf("❌  \033[31mERROR\033[0m : Gagal connect : %v", err)
        helper.EmptySpace()
        helper.DecorTwo()
        helper.EmptySpace()
        os.Exit(1)
    }
    defer reposity.CloseDb()
    
    getUserName()
}