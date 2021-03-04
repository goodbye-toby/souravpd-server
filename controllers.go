package main

import(
	"sync"
	"encoding/json"
)

type BlogHandler struct{
	sync.Mutex
}

func respondWithError(w http.ResponseWriter , code int , msg string){
	respondWithJSON(w , code , map[string]string{"error":msg})
}
func respondWithJSON(w http.ResponseWriter , code int , data interface{}){
	response, err := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
