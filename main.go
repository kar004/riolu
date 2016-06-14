
package main

import (
	"flag"//implementa la linea de comandos
	"log"//implementa registros simples
	"net/http"//implementa protocolo http
	"text/template"//implementa las plantillas de texto
	"fmt"//con esta libreria se puede imprimir en la terminal
	r "github.com/kar004/gorethink"///driver gorethink
)

var addr = flag.String("addr", ":8083", "http service address")//la variable addr contiene la ruta del puerto del servidor
var homeTempl = template.Must(template.ParseFiles("prueba2.html"))//la variable homeTempl toma el archivo html creado
var conn *r.Session

	var row map[string]interface{}
	func init() {
    var errC error
    conn, errC = r.Connect(r.ConnectOpts{
        Address:  "127.0.0.1:28015",
        Database: "test",
        Username: "avsys",
        Password: "sysadmin121",
    })
    if errC != nil {
        fmt.Println("error en la conexion")
        return
    } else {
        fmt.Println("conexion establecida a la DB")
    }
	}

	func serveHome(w http.ResponseWriter, r *http.Request) {//esta funcion determina si esta habilitada la pagina
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")//se otorgan caracteristicas para la respuesta de escritura
	homeTempl.Execute(w, r.Host)////////////////////<------------esto hace que se haga el refresh
	//msg := (string(r.FormValue("id")))
	
	}
	func main() {
	flag.Parse()//analisis gramatico
	go h.run()// se corre el hilo del archivo hub
	http.HandleFunc("/", serveHome)//indica la direccion de ejecucion
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	
}
/////////////////////////////////////////////////////////
	func max() int {
	res, err := r.DB("test").Table("prueba").Max().Run(conn)
		if err != nil {
    	fmt.Println("no hay datos")
    	return 1
		}
	fmt.Println("hay muchos datos")
	var row map[string]interface{}
	res.One(&row)
    max:=row["id"].(float64)
    b := int(max)
    b++
	return b
	}
///////////////////////////////////////////////////////////////
		func insertar (id int, nombre string, edad string) {
		var datos = map[string]interface{}{
			"id": id,
        "nombre":  nombre,
        "edad": edad,
    }
    _, err1 := r.Table("prueba").Insert(datos).RunWrite(conn)
    if err1 != nil {
        fmt.Println("error al insertar")
    }	
}
/////////////////////////////////////////////////////////////
	func eliminar (id int) {
    _, err1 := r.DB("test").Table("prueba").Get(id).Delete().RunWrite(conn)
    if err1 != nil {
        fmt.Println("error al eliminar")
    	}	
	}
////////////////////////////////////////////////////////////
	func modificar (id int, nombre string, edad string) {
		var datos = map[string]interface{}{
        "nombre":  nombre,
        "edad": edad,
    }
    _, err1 := r.Table("prueba").Get(id).Update(datos).RunWrite(conn)
    if err1 != nil {
        fmt.Println("error al modificar")
        return
    }
}
/////////////////////////////////////////////////////////////
	func consulta(id int){
	res, err := r.DB("test").Table("prueba").Get(id).Run(conn)
	if err != nil {
    fmt.Println("no hay datos")
	}
	fmt.Println("hay muchos datos")
	res.One(&row)
	if res.Err() != nil {
	}
	defer res.Close() 
	}