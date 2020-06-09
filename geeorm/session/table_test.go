package session

import (
	"fmt"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	id, _ := s.CreateTable()
	fmt.Println("id", id)
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
