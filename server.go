package main
//imports
import(
	"net/http"
	"log"
	"sync"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"errors"
	"github.com/joho/godotenv"
	"os"
)

//Utility Functions
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

func slugFromURL(r *http.Request) (string , error){
	parts := strings.Split(r.URL.String() , "/")
	if len(parts) != 3{
		return "~",errors.New("Slug Doesnot Exist")
	}	
	slug :=  parts[len(parts)-1]
	return slug,nil
}

func connectDB() *sql.DB{
	DB_STRING,ok := os.LookupEnv("DB_STRING")
	if !ok{
		panic("DB details not found")
	}
	db, err := sql.Open("mysql", DB_STRING)
	if err != nil{
		panic(err)
	}
	return db
}
//Models
type Post struct{
	Id int `json:"id"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	Content string `json:"content"`
	Date string `json:"date"`
}

type BlogHandler struct{
	sync.Mutex
	db *sql.DB
}

func newBlogHandler(db *sql.DB) *BlogHandler{
	return &BlogHandler{
		db:db,
	}
}
//controllers
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
	defer bh.Unlock()
	bh.Lock()
	slug , err := slugFromURL(r)
	if err != nil{
		rows,err := bh.db.Query("SELECT * FROM posts ORDER BY date DESC")
		if err != nil{
			respondWithError(w , http.StatusInternalServerError , "Error in Query")
		} 
		var data []Post
		for rows.Next(){
			var post Post
			err := rows.Scan(&post.Id , &post.Slug , &post.Description , &post.Content , &post.Date)
			if err != nil{
				respondWithError(w , http.StatusInternalServerError , "Error in Query")
			}
			data = append(data , post)
		}
		respondWithJSON(w , http.StatusOK , data)
		defer rows.Close()
		return
	}else{
		data,err := bh.getSingle(w , slug)
		if err != nil{
			respondWithError(w , http.StatusNotFound , "Requested Resource does not Exist")
		}else{
			respondWithJSON(w , http.StatusOK , data)
		}
		return
	}
}

func (bh *BlogHandler) getSingle (w http.ResponseWriter ,slug string) (interface{} , error){
	rows,err := bh.db.Query("SELECT * FROM posts WHERE slug=? LIMIT 1" , slug)
	defer rows.Close()
	if err != nil{
		respondWithError(w , http.StatusInternalServerError , "Error in Query")
	}
	var l int = 0
	var post Post
	for rows.Next(){
		l += 1
		err := rows.Scan(&post.Id , &post.Slug , &post.Description , &post.Content , &post.Date)
		if err != nil{
			respondWithError(w , http.StatusInternalServerError , "Internal Server Error")
		}
	}	
	if l == 0{
		return "~",errors.New("Requested Resource doesnot Exist")
	}else{
		return post,nil		
	}	
}

func (bh *BlogHandler) post(w http.ResponseWriter , r *http.Request){
}
func (bh *BlogHandler) put(w http.ResponseWriter , r *http.Request){
}
func (bh *BlogHandler) delete(w http.ResponseWriter , r *http.Request){
}

func init(){
	err := godotenv.Load(".env")
	if err != nil{
		panic("Error Loading .env Config")
	}
}

func main(){
	
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("Undefined Environment Variables")
	}
	db := connectDB()
	defer db.Close()
	bh := newBlogHandler(db)
	http.Handle("/blog" , bh)
	http.Handle("/blog/" , bh)
	log.Fatal(http.ListenAndServe(port , nil))
}
