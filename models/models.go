package models;

type ShortUrl struct {
    Id string
    OriginalUrl string
    Key string
}

type User struct {
    Id string
    Name string
    Password string
    Email string
}
