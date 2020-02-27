/*
#   cipherMe
#
#------------------------------------------------------
#
#   Application is a tool for education and lear
#   how most popular ciphers work.
#
#------------------------------------------------------
#
#   Author: Dmitry Alekseev <i.mark.vortman@gmail.com>
#
#------------------------------------------------------
*/
package main

import (
	"net/http"
	"log"
	"encoding/json"
    "io/ioutil"
    "strconv"
    "strings"
    "math"
)

/*
#   Default structs used in application
#   for getting request and receiving response
*/
type ErrorMessage struct {
    Title 		string `json:"title"`
    Description string `json:"description"`
}

type TurnRequestData struct {
	Cipher  	string `json:"cipher"`
    Alphabet    string `json:"alphabet"`
    Key         string `json:"key"`
    Text        string `json:"text"`
}

type TurnResponseData struct {
	// Action 		string `json:"action"`
	// Original  	string `json:"original"`
    Converted   string `json:"converted"`
}

const x = "lol"

/*
#   Point of entry, server initialization
*/
func main() {
    http.HandleFunc("/encrypt", turn)
    http.HandleFunc("/decrypt", turn)

    log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
#   "Turn" function is a 'intermediary' layer between client and encryption/decryption
*/
func turn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")


    /*
    #   Block for handling wrong HTTP
    #   method
    */
    if r.Method != "TURN" {
		w.WriteHeader(http.StatusBadRequest)

        invalidMethod := ErrorMessage {
        	Title: 			"Connection Error",
            Description: 	"Can't connect: wrong HTTP method.",
        }

        json.NewEncoder(w).Encode(invalidMethod)
        return
    }

    /*
    #   Formation data for delivery to cryptor/decryptor
    */
	var reqData TurnRequestData
    var action string = r.URL.Path[1:]
    var converted string;

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

    if err != nil {
		http.Error(w, err.Error(), 500)
		return
    }

	err = json.Unmarshal(body, &reqData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

    switch reqData.Cipher {
    case "caesar":
        converted = caesar(reqData.Text, reqData.Key, action)
    case "affine":
        converted = affine(reqData.Text, reqData.Key, action)
    }

    resData := TurnResponseData {
    // Action: 		action,
    // Original: 	reqData.Cipher,
        Converted:      converted,
    }

    json.NewEncoder(w).Encode(resData)
}

func caesar(text string, key string, action string) string {
    shift, err := strconv.Atoi(key)
    if err != nil {
        return err.Error()
    }

    if action == "decrypt" {
        shift = shift * -1;
    }

    shift = (shift%26 + 26) % 26
    result := make([]byte, len(text))
    for i := 0; i < len(text); i++ {
        t := text[i]
        var a int = 'a'
        switch {
        case 'a' <= t && t <= 'z':
            a = 'a'
        case 'A' <= t && t <= 'Z':
            a = 'A'
        default:
            result[i] = t
            continue
        }

        result[i] = byte(a + ((int(t)-a)+shift)%26)
    }
    return string(result)
}

func affine(text string, key string, action string) string {
    keys := strings.Split(key, ",")

    k, err1 := strconv.Atoi(keys[0])
    if err1 != nil {
        return "Error: can't convert first key"
    }

    p, err2 := strconv.Atoi(keys[1])
    if err2 != nil {
        return "Error: can't convert second key"
    }

    if gcd(26, k) != 1 {
        return "Error: wrong key"
    }

    result := make([]byte, len(text))
    for i := 0; i < len(text); i++ {
        t := text[i]
        var a int = 'a'
        switch {
            case 'a' <= t && t <= 'z':
                a = 'a'
            case 'A' <= t && t <= 'Z':
                a = 'A'
            default:
                result[i] = t
                continue
        }

        x := int(t) - a
    
        if action == "encrypt" {
            result[i] = byte(a + ((k * x) + p) % 26)
        } else {
            result[i] = byte(a + (int( math.Pow(float64(k),11) * float64(x - p)) % 26 + 26) % 26 )
        }
    }
    
    return string(result)
}

func gcd(a int, b int) int {
    for a != b {
        if a > b {
            a -= b
        } else {
            b -= a
        }
    }

    return a
}