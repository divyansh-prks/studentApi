package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divyansh/students-api/internal/config"
	"github.com/divyansh/students-api/internal/http/handlers/student"
	"github.com/divyansh/students-api/internal/storage/sqlite"
)

func main(){

	//load config 

	cfg := config.MustLoad()
	//database setup 


	storage , err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage Initialized" , slog.String("env" , cfg.Env), slog.String("version" , "1.0.0"))
	//setup router 
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students" , student.New(storage))
	router.HandleFunc("GET /api/students/{id}" , student.GetById(storage))
	//setup server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}


    slog.Info("server started" , slog.String("address" , cfg.Addr))
	fmt.Printf("Server is started %s" , cfg.HTTPServer.Addr)

	done := make(chan os.Signal , 1)

	signal.Notify(done , os.Interrupt , syscall.SIGINT , syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()

	if err != nil {
		log.Fatal("failed to start the server")
	}


	}()

	<-done
	
	slog.Info("Shutting down the server")

	ctx , cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()
	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown the server" , slog.String("error" , err.Error()))
	}



} 