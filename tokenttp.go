package main

import (
	"net/http"
	"io"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"math/rand"
	"encoding/base64"
	"crypto/hmac"
	"strings"
)

type FromJson struct{
	Token string `json:"token"`
	Name string `json:"name"`
	Pass string `json:"pass"`
}
type User struct {
	Username string
	Passwordhash string
}
type Users []User
var Userslist Users
type Tokens [] string
var Tokenlist Tokens

type Header struct {
	Alg string
	Typ string
}

type Playload struct {
	Name string `json:"Name"`
	Sub int `json:"Sub"`
	Exp int64 `json:"Exp"`
}

func (u Users) CheckUser(u2 User) bool {
	for i := range u {
		if u[i].Username==u2.Username && u[i].Passwordhash==u2.Passwordhash {
			return true
		}
	}
	return false
}

func (t Tokens) CheckToken (token string) bool {
	for i := range t {
		if t[i]== token {
			return true
		}
	}
	return false
}

func Deletetoken (t string) Tokens {
	for i := 0; i<=len(t); i++ {
		if Tokenlist[i]==t {
			Backtokenlist:= append(Tokenlist[:i], Tokenlist[i+1:]...)
			return Backtokenlist
		}
	}
	return nil
}

func GetTime (token string) int64 {
	i:= strings.Index(token, ".")
	token2:=  strings.Replace(token, ".", "-", 1)
	i1:= strings.Index(token2, ".")
	var k string
	k=""
	for i2:=i+1; i2<i1; i2++ {
		k+=string(token[i2])
	}
   DecodedStrPl, err:= base64.StdEncoding.DecodeString(k)
	if err!=nil{
		panic (err)
	}
	var FromJsonPl Playload
	err = json.Unmarshal(DecodedStrPl, &FromJsonPl)
	time:= FromJsonPl.Exp
	return time
}

func Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,PUT,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Lenght, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if req.Method == "POST" {
			data, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Body.Close()
			
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Lenght, Authorization")
			
			var u FromJson 
			err = json.Unmarshal(data, &u)
			
			var IfTokenIsInList bool
			IfTokenIsInList = false
			if len(Tokenlist)!= 0 && u.Token != ""{
				IfTokenIsInList= Tokenlist.CheckToken(u.Token)
			}
			if IfTokenIsInList == true {
				timenow := time.Now()
				timeexp:= GetTime(u.Token)
				timenowint:= timenow.Unix()
				if timenowint>timeexp {
					Tokenlist = Deletetoken(u.Token)
				} else {
					io.WriteString(w, "You got access to the data")
				}
			}else {
				io.WriteString(w, "You shall not pass")
			}
		
		}
}


func Handler2(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,PUT,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Lenght, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Body.Close()
		
		
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Lenght, Authorization")
		
		
		var UserJson FromJson
		err = json.Unmarshal(data, &UserJson)
		if err != nil {
			fmt.Println(err)
			return
		}
		var UserCheck User
		UserCheck.Username = UserJson.Name
		h := sha256.New()
		h.Write([]byte(UserJson.Pass))
		UserCheck.Passwordhash = hex.EncodeToString (h.Sum(nil))
		IsHereSuchUser := Userslist.CheckUser(UserCheck)
		if IsHereSuchUser == true {
				
			var headerr Header
			headerr.Alg = "HMAC"
			headerr.Typ = "JWT"
			var playl Playload
			playl.Name = UserCheck.Username
			playl.Sub = rand.Intn(10000)
				startt := time.Now()
				exptime:=startt.Add(time.Minute*10)
			playl.Exp =exptime.Unix()
			h, err := json.Marshal(headerr)
			if err != nil {
				fmt.Println(err)
				return
			}
			p,err := json.Marshal(playl)
			if err != nil {
			fmt.Println(err)
			return
			}
			var IsHereSuchToken bool
			IsHereSuchToken = false
			if len (Tokenlist)!= 0 { IsHereSuchToken = Tokenlist.CheckToken (UserJson.Token) }
			if IsHereSuchToken == false || len (Tokenlist)== 0 {
				secret:=string(rand.Intn(10000))
				ha := hmac.New(sha256.New, []byte(secret))
				sum := (base64.StdEncoding.EncodeToString(h) + "." + base64.StdEncoding.EncodeToString(p))
				ha.Write([]byte(sum))
				sha := hex.EncodeToString(ha.Sum(nil))
				token := (base64.StdEncoding.EncodeToString(h) + "." + base64.StdEncoding.EncodeToString(p) + "." + base64.StdEncoding.EncodeToString([]byte(sha)))
				
				
				
				Tokenlist = append(Tokenlist,token)
				io.WriteString(w, token)
			} else {
				io.WriteString(w, UserJson.Token)
			}
		}else {
			io.WriteString(w, "No acess")
		} 
			
		}
	
}


func main() {
	var user User 
	user.Username= "a"
	h := sha256.New()
	h.Write([]byte("a"))
	user.Passwordhash= hex.EncodeToString (h.Sum(nil))
	Userslist = append(Userslist , user)
	
	http.HandleFunc("/data", Handler)
	http.HandleFunc("/authorization", Handler2)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
