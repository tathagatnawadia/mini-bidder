package main

import(
    "net/http"
    "encoding/json"
    "time"
    "math/rand"
)

type AdPlacement struct {
    Id string
}

type AdObject struct {
    AdId string
    BidPrice int
    CreatedAt time.Time
}

func main(){
    mux := http.NewServeMux()
    mux.HandleFunc("/bid", bidHandler)
    http.ListenAndServe(":5000", mux)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func isOkay() bool {
    return rand.Float32() < 0.5
}

func RandomCompanyName(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func bidHandler(w http.ResponseWriter, r *http.Request){

    if isOkay() == false {
        w.WriteHeader(204)
    } else {
        adPlacement := AdPlacement{}
        adObject := AdObject{}

        err := json.NewDecoder(r.Body).Decode(&adPlacement)
        if err != nil{
            panic(err)
        }

        adObject.CreatedAt = time.Now().Local()
        adObject.AdId = "COMPANY_" + RandomCompanyName(2) + "_"  +RandomCompanyName(2)
        adObject.BidPrice = rand.Intn(100000)

        adJson, err := json.Marshal(adObject)
        if err != nil{
            panic(err)
        }

        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(adJson)
    }
    
}