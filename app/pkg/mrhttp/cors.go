package mrhttp

import (
    "net/http"

    "github.com/rs/cors"
)

// go get -u github.com/rs/cors

type Cors struct {
    cors *cors.Cors
}

type CorsOptions struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
    ExposedHeaders []string
    AllowCredentials bool
    Debug bool
}

func NewCors(opt CorsOptions) *Cors {
    return &Cors{
        cors: cors.New(cors.Options{
            AllowedOrigins:   opt.AllowedOrigins,
            AllowedMethods:   opt.AllowedMethods,
            AllowedHeaders:   opt.ExposedHeaders,
            ExposedHeaders:   opt.ExposedHeaders,
            AllowCredentials: opt.AllowCredentials,
            Debug:            opt.Debug,
        })}
}

func (c *Cors) Middleware(next http.Handler) http.Handler {
    return c.cors.Handler(next)
}
