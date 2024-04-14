package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func main() {

	var err error
	bot, err = tgbotapi.NewBotAPI("6740393670:AAHTCTwZ2wQQ9LL-YNTjfsBZ3KH38-Q4Los")
	if err != nil {
		log.Panic(err)
	}
	// Устанавливаем обработчики маршрутов
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/order", orderHandler) // Новый обработчик для маршрута /order
	http.HandleFunc("/submit_order", submit)
	http.HandleFunc("/confirmation.html", confirmationHandler)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("HTML/html/css"))))
	http.Handle("/newcss/", http.StripPrefix("/newcss/", http.FileServer(http.Dir("HTML/order/newcss"))))
	http.Handle("/confirmation/", http.StripPrefix("/confirmation/", http.FileServer(http.Dir("HTML/order/confirmation.html"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("HTML/html/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("HTML/html/fonts"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("HTML/html/img"))))
	http.Handle("/2img/", http.StripPrefix("/2img/", http.FileServer(http.Dir("HTML/order/2img"))))

	log.Println("Сервер запущен на порту 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// Обработчик для главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/html/index.html")
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/order/order.html")
}

func submit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	zakaz := r.FormValue("zakaz")
	adres := r.FormValue("adres")

	db, err := sql.Open("mysql", "root:@tcp(database:3306)/china")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `orders` (`name`,`zakaz`,`adres`) VALUES('%s', '%s', '%s')", name, zakaz, adres))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	orderText := fmt.Sprintf("**_Новый заказ!_**\n*Заказчик:* %s\n*Заказ:* %s\n*Адрес:* %s", name, zakaz, adres)

	// Отправляем сообщение в телеграм
	msg := tgbotapi.NewMessage(-1001596545488, orderText)
	msg.ParseMode = "Markdown" // Устанавливаем режим разметки Markdown
	_, err = bot.Send(msg)
	if err != nil {
		panic(err)
	}

	// Перенаправляем пользователя на страницу confirmation.html
	http.Redirect(w, r, "/confirmation.html", http.StatusSeeOther)
}
func confirmationHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/order/confirmation.html")
}
