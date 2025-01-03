package main
import (
    //"fmt"
    "net/http"
)

func main() {
    servemux := http.ServeMux{}
    server := http.Server{
        Handler: &servemux,
        Addr: ":8080",
    }

    servemux.Handle("/", http.FileServer(http.Dir(".")))
    servemux.Handle("/assets/logo.png", http.FileServer(http.Dir(".")))

    server.ListenAndServe()
}


