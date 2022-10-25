# MailSearch-Back

### Tecnologías usadas en el proyecto
* ZincSearch **base de datos**
* GO **lenguaje de programacion**
  * chi router **API Router**


---

Esta parte del proyecto consta de dos funcionalidades principales. Indexar Información en la base de datos ya mencionada y recuperar información de la misma.


La Api se aloja en el puerto :5000, con la dirección  **/api/indexer/** y solo tiene dos metodos
  
  * **GET**. Con este metodo cumple la funcion de mandar la peticion a la base de datos con los parametros asignados correspondientes a la palabra a buscar y a la cantidad de resultados maximos, dichos parametros se le psar por *URL*.


  
      `/api/indexer/{term}-{max}`
 
     
     
  * **POST**. Este es el método con el que indexamos información a la base de datos, que realmente lo único que hace es iniciar el proceso de indexación, el hecho de hacerlo a través del método post de la API es únicamente como mecanismo de llamada a este proceso.


```go
  router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Default().Handler)

	router.Post("/api/indexer/", handlers.Indexar)

	router.Get("/api/indexer/{term}-{max}", handlers.GetCorreos)
	
	http.ListenAndServe(":5000", router)
 ```

---

Primero debemos ingresar información a la base de datos, aquí surge el primer problema ya que dicha información se encuentra descrita en archivos de texto que hacen parte de una cadena de directorios, que más bien se asemejan a un árbol donde partimos de un directorio con más directorios dentro hasta llegar a los archivos de texto con la información requerida.

Como primer paso debemos acceder a cada uno de estos archivos de texto, para ello debemos navegar entre las carpetas. Esto lo hacemos con ayuda del paquete **ioutil** utilizado en esta función recursiva que recorre todas las carpetas abriendo cada uno de los archivos de texto.

```go
  func abrirArchivos(phat string) {
	archivos, _ := ioutil.ReadDir(phat)
	for _, archivo := range archivos {
		phat2 := phat + "\\" + archivo.Name()

		if !archivo.IsDir() {

			Lineas := leerLinea(phat2)

			Co := models.Transformar_Correo(Lineas)

			st, _ := json.Marshal(Co)

			agregarEnBD(string(st))

		} else {
			abrirArchivos(phat2)
		}

	}
}
 ```
 
cada vez que encontramos un archivo de texto leemos su información línea a línea y vamos asignando dicha información con ayuda de la estructura Correo, que contiene las variables correspondientes a la información que nos interesa.


```go
  type Correo struct {
	ID         int64  `json:"id"`
	Message_id string `json:"message_id"`
	Date       string `json:"date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}
 ```
 
 
ya de esta forma es más fácil con ayuda del paquete **json** pasar de este objeto a un Sting con formato json que represente la información del correo respectivamente, ya con este string podemos pasar la información a base de datos haciendo una solicitud a la base de datos con este string

---

Para obtener el resultado de la busqueda es mas faicil, unicamente requerimos capturar dos parametros de la   *url* del metodo **GET**, estos como ya explique antes hacen referencia a lo que vamos a buscar y a la cantidad maxima de resultados, conestos parametros construimos el siquente query.

```go
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
 ```
 
 Este es el query que mandaremos en la peticion a la base de datos, y la respuesta que recibimos se manda como respuesta de nuestro metodo.
