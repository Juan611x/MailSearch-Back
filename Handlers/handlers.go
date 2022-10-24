package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	models "indexer/Models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

type CorreoResponse struct {
	ID         int64  `json:"id"`
	Message_id string `json:"message_id"`
	Date       string `json:"date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}

func GetCorreos(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	max := chi.URLParam(r, "max")

	query := fmt.Sprintf(`{
			"search_type": "matchphrase",
			"query":
			{
				"term": "%s"
			},
			"from": 0,
			"max_results": %s,
			"_source": []
		}`, term, max)

	req, err := http.NewRequest("POST", "http://localhost:4080/api/Prueba/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	hits := string(body)

	hits = strings.Split(hits, ":[")[1]
	hits = strings.Split(hits, "]}")[0]

	fmt.Println(hits)
	fmt.Fprintln(w, string(body))
}

func Indexar(w http.ResponseWriter, r *http.Request) {
	phat := "C:\\Users\\juan6\\OneDrive\\Escritorio\\MailSearch\\enron_mail_20110402\\maildir"
	// phat := "C:\\Users\\juan6\\OneDrive\\Escritorio\\Prueba\\Indexer\\ejemplo"
	abrirArchivos(phat)
}

func abrirArchivos(phat string) {
	archivos, _ := ioutil.ReadDir(phat)
	for _, archivo := range archivos {
		phat2 := phat + "\\" + archivo.Name()

		if !archivo.IsDir() {

			if true {

				Lineas := leerLinea(phat2)

				Co := models.Transformar_Correo(Lineas)

				st, _ := json.Marshal(Co)

				// fmt.Println(string(st))

				agregarEnBD(string(st))

			} else {
				fmt.Println(phat2)
			}

		} else {
			abrirArchivos(phat2)
		}

	}
}

func leerLinea(fileName string) []string {
	var lineas []string

	file, error := os.Open(fileName)
	if error != nil {
		fmt.Println(error.Error())
	}

	defer file.Close()

	contenido := bufio.NewScanner(file)

	for contenido.Scan() {
		lineas = append(lineas, contenido.Text())
	}

	return lineas
}

func agregarEnBD(data string) {
	req, err := http.NewRequest("POST", "http://localhost:4080/api/Prueba/_doc", strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
