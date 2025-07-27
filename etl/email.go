package main

import "github.com/jdetok/golib/maild"

func EmailLog(file string) error {
	m := maild.MakeMail(
		[]string{"jdekock17@gmail.com"},
		"Go bball ETL log attached",
		"the Go bball ETL process ran. The log is attached.",
	)
	return m.SendMIMEEmail(file)
}
