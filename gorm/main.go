package main

import "fmt"

func main() {
	fmt.Println(getUsers(&User{Name: "meow", Age: 0}))
	fmt.Println("hello world")
}

func getUsers(opt *User) (users []User, err error) {
	err = db.Model(&User{}).
		Where(opt).
		Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	return users, nil
}
