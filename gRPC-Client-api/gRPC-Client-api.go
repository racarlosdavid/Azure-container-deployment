package main

import (
	"context"
	"log"
	"time"
	"net/http"
	"encoding/json" 
	"github.com/gorilla/mux" 
	"os"
	"github.com/Pallinder/go-randomdata"
	"google.golang.org/grpc"
	pb "github.com/racarlosdavid/demo-gRPC-kubernetes/gRPC-Client-api/proto"
)

var nombre_api = "default"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API GO - gRPC Client!\n"))
}

func operacionesAritmeticas(w http.ResponseWriter, r *http.Request) {
	operacion := mux.Vars(r)["op"] //Obtengo la operacion a realizar
	num1 := mux.Vars(r)["val1"] //Obtengo el valor 1
	num2 := mux.Vars(r)["val2"] //Obtengo el valor 2
	host:= os.Getenv("HOST")
	/********************************** gRPC llamada al servidor ********************************/
	conn, err := grpc.Dial(host+":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOperacionAritmeticaClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := c.OperarValores(ctx, &pb.OperacionRequest{
		Operacion: operacion,
		Valor1:num1,
		Valor2:num2,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	//log.Printf("Greeting: %s", reply.GetResultado())
	/********************************** gRPC ********************************/

	/********************************** Respuesta ********************************/
	w.Header().Set("Content-Type", "application/json")
   	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(struct {
		Mensaje string `json:"mensaje"`
		Server string `json:"server"`
	}{Mensaje: reply.GetResultado(),Server:nombre_api})
}

func main() {
	nombre_api = randomdata.SillyName()
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/operacion/{op}/valor1/{val1}/valor2/{val2}",operacionesAritmeticas).Methods("GET")
    log.Println("Listening at port 2000") 
	log.Fatal(http.ListenAndServe(":2000", router))
}