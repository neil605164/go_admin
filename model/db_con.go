package model

import (
	"GO_Admin/global"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// dbConnect 建立 DB 連線
func dbConnect() (db *gorm.DB, err error) {
	USER := global.Config.Database.User
	PASSWORD := global.Config.Database.Password
	HOST := global.Config.Database.Host
	DATABASE := global.Config.Database.Database

	// 組合連線資訊
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASSWORD, HOST, DATABASE)

	// 建立連線
	db, err = gorm.Open("mysql", connectionString)
	if err != nil {
		err = global.NewError{
			Title:   "DB connect Fail",
			Message: fmt.Sprintf("Error message is: %s", err),
		}
		return nil, err
	}

	db.LogMode(true)
	return db, nil
}

// CheckTableIsExist 啟動main.go服務時，直接檢查所有 DB 的 Table 是否已經存在
func CheckTableIsExist() error {
	db, err := dbConnect()
	if err != nil {
		return err
	}

	defer db.Close()

	if !db.HasTable("users") {
		db.AutoMigrate(&User{})
	}

	if !db.HasTable("user_infos") {
		db.AutoMigrate(&UserInfo{})
	}

	if !db.HasTable("files") {
		db.AutoMigrate(&File{})
	}

	return nil
}

// SQLRegisterMem 註冊會員
func SQLRegisterMem(rgMem *global.RegisterMemberOption) (err error) {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	user := User{
		Username: rgMem.Username,
		Password: rgMem.Password,
	}

	userInfo := UserInfo{
		Username: rgMem.Username,
		Nickname: rgMem.Nickname,
		Email:    rgMem.Enail,
		Addr:     rgMem.Addr,
	}

	if CheckMemExist(rgMem.Username, db); err != nil {
		return err
	}

	tx := db.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when register user",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}

	if err := tx.Create(&userInfo).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when register user",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}

	tx.Commit()
	return err
}

// SQLGetUserList 取得用戶清單
func SQLGetUserList() (userList *[]User, err error) {
	var users []User

	db, err := dbConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := db.Find(&users).Error; err != nil {
		err = global.NewError{
			Title:   "Unexpected error when get all users list",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return nil, err
	}
	fmt.Println(&users)
	return &users, nil
}

// SQLEditUserInfo 編輯會員資料
func SQLEditUserInfo(edUserInfo *global.EditUserInfoOption) (err error) {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	user := User{
		Username: edUserInfo.Username,
		Password: edUserInfo.Password,
	}

	userInfo := UserInfo{
		Username: edUserInfo.Username,
		Nickname: edUserInfo.Nickname,
		Email:    edUserInfo.Enail,
		Addr:     edUserInfo.Addr,
	}

	tx := db.Begin()

	if CheckMemExist(edUserInfo.Username, db); err != nil {
		return err
	}

	if err := tx.Model(user).Where("username = ?", user.Username).Updates(user).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when edit users table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}

	if err := tx.Model(userInfo).Where("username = ?", userInfo.Username).Updates(userInfo).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when edit user_infos table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}

	tx.Commit()
	return nil
}

// SQKFreezeUserAccount 停用用戶帳號
func SQKFreezeUserAccount(freezeMem *global.FreezeUserAccountOption) (err error) {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	user := User{
		Status: 1,
	}

	if CheckMemExist(freezeMem.Username, db); err != nil {
		return err
	}

	tx := db.Begin()

	if err = tx.Model(user).Where("username = ?", freezeMem.Username).Update(user).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when edit users table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}
	tx.Commit()

	return nil
}

// SQLDeleteUserAccount 刪除用戶帳號
func SQLDeleteUserAccount(deleteMem *global.DeleteUserAccountOption) (err error) {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	user := User{}
	userInfo := UserInfo{}

	// 可以移除，因為下方有檢查『影響數量』
	if CheckMemExist(deleteMem.Username, db); err != nil {
		return err
	}

	tx := db.Begin()
	execRes := tx.Model(user).Where("users.username = ?", deleteMem.Username).Delete(user)
	if execRes.Error != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when delete users table",
			Message: fmt.Sprintf("Error massage is: %s", execRes.Error),
		}
		return err
	}

	execRes = tx.Model(userInfo).Where("username = ?", deleteMem.Username).Delete(userInfo)
	if execRes.Error != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when delete user_infos table",
			Message: fmt.Sprintf("Error massage is: %s", execRes.Error),
		}
		return err
	}

	tx.Commit()
	return nil
}

// SQLEnableUserAccount 啟用會員帳號
func SQLEnableUserAccount(enableMem *global.EnableUserAccountOption) (err error) {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	users := User{}

	if CheckMemExist(enableMem.Username, db); err != nil {
		return err
	}

	tx := db.Begin()

	if err = tx.Model(users).Where("username = ?", enableMem.Username).Update("status", 0).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when edit users table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}
	tx.Commit()

	return nil
}

// SQLUploadFile 上傳檔案
func SQLUploadFile(fileInfo *global.UploadFileOption) error {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	file := File{
		FileName: fileInfo.FileName,
		FileSize: fileInfo.FileSize,
		FilePath: fileInfo.FilePath,
		FileExt:  fileInfo.FileExt,
	}

	tx := db.Begin()

	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when insert file table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}
	tx.Commit()

	return nil
}

// SQLUploadMultiFile 上傳多個檔案
func SQLUploadMultiFile(fileInfo *global.UploadMultiFileOption) error {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	// 組query語法
	var strArr []string
	query := "INSERT INTO files (`created_at`, `file_name`, `file_size`, `file_path`, `file_ext`) VALUES "
	for i := 0; i < len(fileInfo.File); i++ {
		now := time.Now()

		strfileSize := strconv.FormatInt(fileInfo.FileSize[i], 10)
		strArr = append(strArr, "('"+now.Format("2006-01-02 15:04:05")+"','"+fileInfo.FileName[i]+"','"+strfileSize+"','"+fileInfo.FilePath+"','"+fileInfo.FileExt[i]+"')")
	}

	newStr := strings.Join(strArr, ",")
	query += newStr

	tx := db.Begin()
	if err := tx.Exec(query).Error; err != nil {
		tx.Rollback()
		err = global.NewError{
			Title:   "Unexpected error when insert mutli file table",
			Message: fmt.Sprintf("Error massage is: %s", err),
		}
		return err
	}
	// if err := tx.Create(&file).Error; err != nil {
	// tx.Rollback()
	// err = global.NewError{
	// 	Title:   "Unexpected error when insert file table",
	// 	Message: fmt.Sprintf("Error massage is: %s", err),
	// }
	// return err
	// }
	tx.Commit()

	return nil
}
