package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    "github.com/stianeikeland/go-rpio"
)

var (
    listen = ":8082"
)

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: Home")
    fmt.Fprintf(w, "Welcome to the HomePage!")
    // TODO add manual
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/on", turnLedOn).Methods("Get", "Post")
    myRouter.HandleFunc("/off", turnLedOff).Methods("Get", "Post")
    myRouter.HandleFunc("/timer", timer).Methods("Get", "Post")
    myRouter.HandleFunc("/blink", blinkLed).Methods("Get", "Post")
    myRouter.HandleFunc("/status", returnStatus)
    log.Fatal(http.ListenAndServe(listen, myRouter))
}

func turnLedOn(w http.ResponseWriter, r *http.Request) {
    intPin, errorPin := getInt(r, "pin")
    if errorPin != nil {
        fmt.Println(errorPin)
        fmt.Fprintln(w, errorPin)
        return
    }

    if errorRPIO := rpio.Open(); errorRPIO != nil {
        fmt.Println(errorRPIO)
        fmt.Fprintln(w, "...ups, something went horribly wrong. We failed to open a connection to gpio. Further details can be found in the application logs")
        return
    }

    defer rpio.Close()
    gpioPin := rpio.Pin(intPin)

    gpioPin.Output()
    gpioPin.High()

    fmt.Printf("Turning LED on PIN '%d' on\n", intPin)
    fmt.Fprintf(w, "Turning LED on PIN '%d' on\n", intPin)
}

func turnLedOff(w http.ResponseWriter, r *http.Request) {
    intPin, errorPin := getInt(r, "pin")
    if errorPin != nil {
        fmt.Println(errorPin)
        fmt.Fprintln(w, errorPin)
        return
    }

    if errorRPIO := rpio.Open(); errorRPIO != nil {
        fmt.Println(errorRPIO)
        fmt.Fprintln(w, "...ups, something went horribly wrong. We failed to open a connection to gpio. Further details can be found in the application logs")
        return
    }

    defer rpio.Close()
    gpioPin := rpio.Pin(intPin)

    gpioPin.Output()
    gpioPin.Low()

    fmt.Printf("Turning LED on PIN '%d' off\n", intPin)
    fmt.Fprintf(w, "Turning LED on PIN '%d' off\n", intPin)
}

func blinkLed(w http.ResponseWriter, r *http.Request) {
    intPin, errorPin := getInt(r, "pin")
    if errorPin != nil {
        fmt.Println(errorPin)
        fmt.Fprintln(w, errorPin)
        return
    }
    intInterval, errorInterval := getInt(r, "interval")
    if errorInterval != nil {
        fmt.Println(errorInterval)
        fmt.Fprintln(w, errorInterval)
        return
    }

    intCount, errorCount := getInt(r, "count")
    if errorCount != nil {
        fmt.Println(errorCount)
        fmt.Fprintln(w, errorCount)
        return
    }

    if errorRPIO := rpio.Open(); errorRPIO != nil {
        fmt.Println(errorRPIO)
        fmt.Fprintln(w, "...ups, something went horribly wrong. We failed to open a connection to gpio. Further details can be found in the application logs")
        return
    }

    defer rpio.Close()
    gpioPin := rpio.Pin(intPin)

    gpioPin.Output()
    for x, counter := 0, intCount*2; x < counter; x++ {
        gpioPin.Toggle()
        time.Sleep(time.Duration(intInterval) * time.Millisecond)
    }

    gpioPin.Low()
}

func timer(w http.ResponseWriter, r *http.Request) {
    intPin, errorPin := getInt(r, "pin")
    if errorPin != nil {
        fmt.Println(errorPin)
        fmt.Fprintln(w, errorPin)
        return
    }
    intTime, errorTime := getInt(r, "time")
    if errorTime != nil {
        fmt.Println(errorTime)
        fmt.Fprintln(w, errorTime)
        return
    }

    if errorRPIO := rpio.Open(); errorRPIO != nil {
        fmt.Println(errorRPIO)
        fmt.Fprintln(w, "...ups, something went horribly wrong. We failed to open a connection to gpio. Further details can be found in the application logs")
        return
    }

    defer rpio.Close()
    gpioPin := rpio.Pin(intPin)

    gpioPin.Output()
    gpioPin.High()
    time.Sleep(time.Duration(intTime) * time.Second)
    gpioPin.Low()
}

func returnStatus(w http.ResponseWriter, r *http.Request) {
    intPin, errorPin := getInt(r, "pin")
    if errorPin != nil {
        fmt.Println(errorPin)
        fmt.Fprintln(w, errorPin)
        return
    }

    if errorRPIO := rpio.Open(); errorRPIO != nil {
        fmt.Println(errorRPIO)
        fmt.Fprintln(w, "...ups, something went horribly wrong. We failed to open a connection to gpio. Further details can be found in the application logs")
        return
    }

    defer rpio.Close()
    gpioPin := rpio.Pin(intPin)
    gpioPin.Input()

    status := "undefined"
    if gpioPin.Read() == 0 {
        status = "off"
    } else {
        status = "on"
    }
    fmt.Printf("Status for LED on PIN '%d': %s\n", intPin, status)
    fmt.Fprintf(w, "Status for LED on PIN '%d': %s\n", intPin, status)
}

func getInt(r *http.Request, urlParameter string) (int, error) {
    keys, ok := r.URL.Query()[urlParameter]

    if !ok || len(keys[0]) < 1 {
        return 0, errors.New("URL parameter '" + urlParameter + "' is missing")
    }

    paramString := keys[0]
    paramInt, err := strconv.Atoi(keys[0])

    if err != nil || paramInt <= 0 {
        return 0, errors.New("URL parameter '" + urlParameter + "' must be a positive number greater than zero. However we got: '" + paramString + "'")
    }

    return paramInt, nil
}

func main() {
    fmt.Println("Started")
    handleRequests()
}
