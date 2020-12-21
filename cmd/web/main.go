package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	debug        bool
	errorLog     *log.Logger
	infoLog      *log.Logger
	accessKey    string
	googleClient *admin.Service
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug stack traces shown to users")
	addr := flag.String("addr", ":4001", "HTTP Network Address")
	credentialsPath := flag.String("credentialsPath", "/secrets/credentials.json", "Path to credentials file")
	accesskey := flag.String("accessKey", "servente-secret-access-key-001!", "Password to access the API")
	tlsCertPath := flag.String("tls-cert-path", "./tls/cert.pem", "TLS certificate path")
	tlsKeyPath := flag.String("tls-key-path", "./tls/key.pem", "TLS key path")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	b, err := ioutil.ReadFile(*credentialsPath)
	if err != nil {
		errorLog.Fatalf("credentials: Unable to read credentials JSON (%s) %v", *credentialsPath, err)
	}
	config, err := google.ConfigFromJSON(b, admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		errorLog.Fatalf("credentials: unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, errorLog)
	ctx := context.Background()
	adminService, err := admin.NewService(ctx, option.WithHTTPClient(client), option.WithScopes(admin.AdminDirectoryGroupReadonlyScope))
	if err != nil {
		errorLog.Fatalf("credentials: unable retrieve directory: %v", err)
	}

	app := &application{
		debug:        *debug,
		errorLog:     errorLog,
		infoLog:      infoLog,
		accessKey:    *accesskey,
		googleClient: adminService,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*tlsCertPath, *tlsKeyPath)
	errorLog.Fatal(err)
}
