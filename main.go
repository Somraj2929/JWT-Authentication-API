package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

//storing all users
var userMap = map[string]string{
	"somraj":"password",
	"arun":"pass1",
	"radhe":"krishna",
}


//Declaring User data type with username and pass
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case "POST":
			var user User

			//decoding the request body into the structs, if error then 404
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil{
				//fmt.Println(user)
				fmt.Fprintf(w, "Kya kar rha tu yeh")
				return
			}
			

			//checking users form map if not match return error
			if userMap[user.Username] == "" || userMap[user.Username] != user.Password{
				fmt.Fprintf(w, "Bhai kisko le aaya tu yeh allowed nhi hai...")
				return
			}

			//if user found in map then generateJWT function will call and
			//token will be generated for user
			token, err := generateJWT(user.Username)
			if err != nil{
				fmt.Fprintf(w, "Token Shi se generate kar")
			}

			//if token successfully generated then printing token
			fmt.Fprintf(w, token)



		case "GET":
			fmt.Fprintf(w, "Sorry! We don't allow you, Please come with GET Method")
			return
	}
}


//generating JWT tokens

var sampleSecretKey = []byte("bhaiyehtosamplehai")

func generateJWT(username string) (string, error) {
	//generating new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"]= true
	claims["username"] = username
	//setting expiry of claims
	claims["exp"] = time.Now().Add(time.Minute* 1 ).Unix()

	//return a complete signed JWT token for signing in.
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		fmt.Errorf("Something went really wrong", err.Error())
		return "", err
	}

	return tokenString, nil

}

//now verifing out generated token

func validateToken(w http.ResponseWriter, r *http.Request)(err error){

	if r.Header["Token"]==nil {
		fmt.Fprintf(w, "Can't find token in header")
		return errors.New("token error")
	}
	
	//parsing token in jwt singingmethodHMAC 
	token, err := jwt.Parse(r.Header["Token"][0], func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing",err)
		}
		return sampleSecretKey, nil
	})

	//if token not present
	if token == nil {
		fmt.Fprintf(w, "Invalid token")
		return errors.New("token error")

	}
	
	// Checking claims of token whether expired or not?

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "Not able to parse claims")
		return errors.New("token error")
	}


	//checking token validity beacuse when we set expiry time
	//then token will be expired
	
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix(){
		fmt.Fprintf(w, "Bhai tera token upar khisk liya")
		return errors.New("token error")
	}

	return nil

}


//getting all books as JSON
func getAllBooksHandler(w http.ResponseWriter, r *http.Request){

	//validating token from both writer and request if both match found then 
	//function will display books data
	err := validateToken(w, r)
	if err == nil{
		w.Header().Set("Content-Type", "application/json")
		books := getAllBook()
		json.NewEncoder(w).Encode(books)
	}
}

//setting some books with authors but here i used movies name with actor😂
func getAllBook() []Book{
	return []Book{
		Book{
			Name: "Biwi No. 1",
			Author: "Govinda",
		},
		Book{
			Name: "Bahubali",
			Author: "Prabhas",
		},
		Book{
			Name: "Hero No. 1",
			Author: "Govinda",
		},
	}
}


//defining book struct 
type Book struct {
	Name string
	Author string
}



func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/getAllBooks", getAllBooksHandler)
	http.ListenAndServe(":8080", nil)
}