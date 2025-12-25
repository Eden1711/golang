package main

import (
	"fmt"
)

// 1. Định nghĩa Hợp đồng: Ai muốn làm Database thì phải biết hàm Save()
type Database interface {
	Save(data string)
}

// 2. MySQL tham gia
type MySQL struct{}

func (m MySQL) Save(data string) {
	fmt.Println("Đã lưu vào MySQL:", data)
}

// 3. FileText tham gia
type FileText struct{}

func (f FileText) Save(data string) {
	fmt.Println("Đã ghi vào file text:", data)
}

// 4. Hàm xử lý logic chính (Không phụ thuộc cụ thể vào MySQL hay File)
// Nó nhận vào Interface
func XuLyDuLieu(db Database, data string) {
	db.Save(data)
}

func main() {
	mySQL := MySQL{}
	file := FileText{}

	// Cắm cái nào vào cũng chạy được!
	XuLyDuLieu(mySQL, "User A")
	XuLyDuLieu(file, "User B")

}
