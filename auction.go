package main

import(
    "fmt"
    "time"
    "github.com/gojektech/heimdall/httpclient"
    "net/http"
    "encoding/json"
    "bytes"
    "sync"
    "math/rand"
)

type AdPlacement struct {
    Id string
}

type AdObject struct {
    AdId string
    BidPrice int
    CreatedAt time.Time
    ElapsedTime float64
}

var mutex = &sync.Mutex{}

func MakeRequest(biddingChannel chan<-AdObject, wg *sync.WaitGroup) {
    defer wg.Done()
    start := time.Now()
    // Use the clients GET method to create and execute the request
    timeout := 180 * time.Millisecond
    values := map[string]string{"Id": "GeeksForGeeksHomePageAd"}
    jsonValue, _ := json.Marshal(values)
    client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))
    resp, err := client.Post("http://localhost:5000/bid", bytes.NewBuffer(jsonValue), nil)
    if err != nil {
        return 
    } else {
        if resp.StatusCode == http.StatusOK{
            decoder := json.NewDecoder(resp.Body)
            var data AdObject
            err = decoder.Decode(&data)
            data.ElapsedTime = time.Since(start).Seconds() * 1000
            biddingChannel <- data
            // Can we do processing here of the highest price bidded ?
        }
    }
    return 
}

func auction(w http.ResponseWriter, r *http.Request){

    adSlot, ok := r.URL.Query()["Id"]
    
    if !ok || len(adSlot[0]) < 1 {
        fmt.Println("Url Param 'Id' is missing")
        return
    }

    var wg sync.WaitGroup
    numberOfRequest := rand.Intn(400)
    wg.Add(numberOfRequest)

    biddingChannel := make(chan AdObject)

    for i := 0; i < numberOfRequest; i++ {
        go MakeRequest(biddingChannel, &wg)
    }
    
    maxBid := 0
    numberOfBids := 0
    go func() {
        for bid := range biddingChannel {
            numberOfBids++
            if maxBid < bid.BidPrice{
                maxBid = bid.BidPrice
            }
        }
    }()

    wg.Wait()
    close(biddingChannel)

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "The winner for " + adSlot[0] +" spot is : %d from %d number of bids\n of %d requested", maxBid, numberOfBids, numberOfRequest)
}


func main(){
    mux := http.NewServeMux()
    mux.HandleFunc("/auction", auction)
    http.ListenAndServe(":6969", mux)
}
