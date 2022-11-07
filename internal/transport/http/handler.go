package http

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Service CompanyRestService
	Server  *http.Server
}

func NewHandler(service CompanyRestService) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(LogMiddlewire)
	h.Router.Use(JSONMiddlewire)

	h.Server = &http.Server{
		Addr:         "0.0.0.0:8888",
		Handler:      h.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return h
}

func (h *Handler) mapRoutes() {
	private := h.Router.PathPrefix("/api/v1/private/company").Subrouter()
	public := h.Router.PathPrefix("/api/v1/public/company").Subrouter()

	h.Router.HandleFunc("/healthcheck", healthCheck)
	public.HandleFunc("/{id}", h.GetCompany).Methods("GET")

	private.HandleFunc("/create", JWTAuth(h.PostCompany)).Methods("POST")
	private.HandleFunc("/delete/{id}", JWTAuth(h.DeleteCompany)).Methods("DELETE")
	private.HandleFunc("/patch/{id}", JWTAuth(h.PartialUpdateCompany)).Methods("PATCH")

}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Fatalf("error serving http%+v\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("Shut down gracefully...")
	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {

	podIp := ""
	ifaces, err := net.Interfaces()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok!",
			"ips":    "nil",
		})
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("could not get address from iface")
			continue
		}
		for _, addr := range addrs {
			ip := addr.(*net.IPNet).IP.String()
			netIp, err := netip.ParseAddr(ip)
			if err != nil {
				log.Println("could not parse ip")
				continue
			}
			calico_cidr_block, _ := netip.ParsePrefix("192.168.0.0/16")

			if calico_cidr_block.Contains(netIp) {
				podIp = netIp.String()
			}
			// process IP address
		}
	}
	host, _ := os.Hostname()
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok!",
		"podIp":   podIp,
		"podName": host,
	})
}
