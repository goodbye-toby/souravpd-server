package main

import(
	"fmt"
	"net/http"
	"log"
	"sync"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter , code int , msg string){
	respondWithJSON(w , code , map[string]string{"error":msg})
}

func respondWithJSON(w http.ResponseWriter , code int , data interface{}){
	response, err := json.Marshal(data)
	if err != nil{
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type Blog struct{
	Title string `json:"title"`
	Description string `json:"description"`
	Html string `json:"html"`
	CreatedAt string `json:"date"`
}

type BlogHandler struct{
	sync.Mutex
}

func (bh *BlogHandler) ServeHTTP(w http.ResponseWriter , r *http.Request){
	switch r.Method{
		case "GET":
			bh.get(w,r)
		case "POST":
			bh.post(w,r)
		case "PUT","PATCH":
			bh.put(w,r)
		case "DELETE":
			bh.delete(w,r)
		default:
			respondWithError(w , http.StatusMethodNotAllowed , "invalid method")
	}	
}

func (bh *BlogHandler) get(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Hello from get\n");
}
func (bh *BlogHandler) post(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Hello from post\n");
}
func (bh *BlogHandler) put(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w,"Hello from put\n");
}
func (bh *BlogHandler) delete(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w, "Hello from delete\n");
}


func main(){
	port := ":8080"
	bh := new(BlogHandler)
	http.Handle("/blog" , bh)
	http.Handle("/blog/" , bh)
	log.Fatal( http.ListenAndServe(port , nil))
}
