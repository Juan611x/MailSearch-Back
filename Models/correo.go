package models

import (
	"strings"
)

type Correo struct {
	ID         int64  `json:"id"`
	Message_id string `json:"message_id"`
	Date       string `json:"date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}

func NewCorreo(Me_I, date, from, to, subject, body string) *Correo {
	P := &Correo{
		Message_id: Me_I,
		Date:       date,
		From:       from,
		To:         to,
		Subject:    subject,
		Body:       body,
	}
	return P
}

func Transformar_Correo(array []string) *Correo {

	me_i := ""
	date := ""
	fron := ""
	to := ""
	subject := ""
	// body := ""

	x := 0

	for i := 0; i < len(array); i++ {
		if strings.Contains(array[i], "Message-ID:") {
			me_i = quitarEncabezado(array[i])
		} else if strings.Contains(array[i], "Date:") {
			date = quitarEncabezado(array[i])
		} else if strings.Contains(array[i], "From:") {
			fron = quitarEncabezado(array[i])
		} else if strings.Contains(array[i], "To:") {
			to = quitarEncabezado(array[i])
		} else if strings.Contains(array[i], "Subject") {
			subject = quitarEncabezado(array[i])
		} else if !strings.Contains(array[i], ":") {
			if i > 14 {
				x = i
				break
			}
		}
	}
	C := NewCorreo(me_i, date, fron, to, subject, agruparBody(array[x:]...))
	return C
}

func quitarEncabezado(texto string) string {

	textoFinal := strings.Split(texto, ":")

	text := ""
	for i := 0; i < len(textoFinal)-1; i++ {
		text += textoFinal[1+i]
	}
	return text
}

func agruparBody(lineas ...string) string {
	st := ""
	for i := 0; i < len(lineas); i++ {
		st += lineas[i] + "\n"
	}

	return st
}
