package main

// import (
//     "time"
// )

type Person struct {
    ID           int            `json:"id"`
    Name         string         `json:"name"`
    Age          int            `json:"age"`
    Motto        string         `json:"motto"`
    LatLocation  float64        `json:"lat"`
    LongLocation float64        `json:"long"`
    Profile      string         `json:"profile"`
    Details      map[string]int `json:"details"`
}

type User struct {
    ID           int    `json:"id"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"`
}

type ProcessedProfile struct {
    ID       int            `json:"id"`
    Name     string         `json:"name"`
    Motto    string         `json:"motto"`
    Distance float64        `json:"distance"`
    Profile  string         `json:"profile"`
    Details  map[string]int `json:"details"`
}

type ProfilePhoto struct {
    Url     string `json:"url"`
    Caption string `json:"caption"`
}

type ChatMessage struct {
    ID      int64  `json:"id"`
    Time    string  `json:"time"`
    Who     string `json:"who"`
    Message string `json:"message"`
}
