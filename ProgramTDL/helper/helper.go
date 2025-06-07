package helper
import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "log"
    "context"
    "ProgramTDL/reposity"
    "ProgramTDL/model"
)

func DecorOne() {
    fmt.Println("°•. ✧ .•°•. ✦ .•°•. ✧ .•°•. ✦ .•°")
}

func DecorTwo() {
    fmt.Println("✧•°: *✧･°:* *:•°✧*:･°✧･°: *✧･°:* *:•°✧")
}

func SpecialDecorOne() {
    fmt.Println(`
🌿 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 🌿
         ● # - To Do List - # ●
🍃 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 🍃
`)
}

func SpecialDecorTwo() {
    fmt.Println(`
╔════════════════════════════════════════╗
║          # Opsi, To Do List :          ║
║════════════════════════════════════════║
║ 1. Tambah Task / Tugas                 ║
║ 2. Lihat  Task / Tugas                 ║
║ 3. Update Task / Tugas                 ║
║ 4. Hapus  Task / Tugas                 ║
║ 5. Keluar / Exit                       ║
╚════════════════════════════════════════╝
`)
}

func SpecialDecorThree() {
	fmt.Println(`
╔════════════════════════════════════════╗
║          # Opsi, Lihat Task :          ║
║════════════════════════════════════════║
║ 1. Lihat semua Tugas                   ║
║ 2. Lihat Tugas berdasarkan filter      ║
╚════════════════════════════════════════╝
`)
}

func SpecialDecorFourht() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║          # Opsi, Lihat Task Berdasarkan Filter :          ║
║═══════════════════════════════════════════════════════════║
║ 1. Lihat Tugas yang Selesai                               ║
║ 2. Lihat Tugas yang Belum Selesai                         ║
║ 3. Lihat Tugas dengan Prioritas Penting                   ║
║ 4. Lihat Tugas dengan Deadline terdekat                   ║
╚═══════════════════════════════════════════════════════════╝
`)
}

func SpecialDecorFive() {
	fmt.Println(`
╔════════════════════════════════════════╗
║          # Opsi, Hapus Task :          ║
║════════════════════════════════════════║
║ 1. Hapus Tugas Soft ( Tidak Permanen ) ║
║ 2. Hapus Tugas Hard ( Permanen )       ║
╚════════════════════════════════════════╝
`)
}

func SpecialDecorSix() {
    fmt.Println(`
╔════════════════════════════════════════╗
║          # Opsi, Update Task :         ║
║════════════════════════════════════════║
║ 1. Update Judul                        ║
║ 2. Update Deskripsi                    ║
║ 3. Update Kategori                     ║ 
║ 4. Update Deadline                     ║ 
║ 5. Update Prioritas                    ║ 
║ 6. Update Status                       ║ 
║ 7. Restore Tugas yang di SoftDelete    ║ 
╚════════════════════════════════════════╝
`)
}

// func

func EmptySpace() {
    fmt.Println("")
}

var reader = bufio.NewReader(os.Stdin)
func InputReader(input string) string {
    fmt.Print(input)
    prompt, _ := reader.ReadString('\n')
    return strings.TrimSpace(prompt)
}

func ReadTaskH() {
    for {
        obI := "id"
        asc := "ASC"
        tableName := "tasks"
        var c model.Task 
        var status string
        var prioritas string
        var no int
        db, err := reposity.GetDb()
    	if err != nil {
    		log.Fatalf("❌  \033[31mERROR\033[0m : %v", err)
    	}
        query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY %s %s ) AS no, id, judul, category, prioritas, status FROM %s WHERE deletesoft = $1", obI, asc, tableName)
        row, err := db.Query(context.Background(), query, false)
    	if err != nil {
    		EmptySpace()
    		log.Printf("\033[31mERROR\033[0m : %v", err)
    		fmt.Printf("❌  \033[31mERROR\033[0m : Gagal mengambil data Tugas : %v", err)
    		EmptySpace()
    		retry := InputReader("Coba lagi? (y/n): ")
    		if strings.ToLower(retry) == "y" {
    			continue
    		} else {
    			EmptySpace()
    			fmt.Println("🚫 Operasi dibatalkan")
    			return
    		}
    	}
    	defer row.Close()
    	EmptySpace()
    	DecorTwo()
    	EmptySpace()
    	fmt.Println("Daftar Tugas : ")
    	EmptySpace()
    	DecorTwo()
    	EmptySpace()
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
    		break
    	}
    	break
    }
}