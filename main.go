package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


var Cfg struct {
	SRV ServerConfig
	DB  DBConfig
}

type ServerConfig struct {
	Secret string
	Prefix string
	Root   string
	Port   string
	Unix   string
	Debug  bool
}

type DBConfig struct {
	DB       string
	User     string
	Password string
}

type Srv struct {
	db *sqlx.DB
}

func main() {
	log.SetFlags(log.Lshortfile)

	config_file := flag.String("c", "config.ini", "Config filename")
	flag.Parse()

	if _, err := toml.DecodeFile(*config_file, &Cfg); err != nil {
		log.Println(err)
		return
	}
	conninfo := fmt.Sprintf("user=%s password=%s dbname=%s  sslmode=disable", Cfg.DB.User, Cfg.DB.Password, Cfg.DB.DB)

	db, err := sqlx.Open("postgres", conninfo)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	s := &Srv{}
	s.db = db
	http.HandleFunc(Cfg.SRV.Prefix+"login", cors(s.login))
	http.HandleFunc(Cfg.SRV.Prefix+"registration", cors(s.registration))
	http.HandleFunc(Cfg.SRV.Prefix+"account", cors(s.account))
	http.HandleFunc(Cfg.SRV.Prefix+"article", cors(s.articleHandler))
	http.HandleFunc(Cfg.SRV.Prefix+"article/comment", cors(s.writeComment))
	http.HandleFunc(Cfg.SRV.Prefix+"article/vote", cors(s.voteArticle))
	http.HandleFunc(Cfg.SRV.Prefix+"article/image", cors(s.deleteImage))
	http.HandleFunc(Cfg.SRV.Prefix+"article/image/all", cors(s.deleteAllImages))

	if Cfg.SRV.Debug {
		log.Printf("Listen: http://localhost:%s\n", Cfg.SRV.Port)
		http.Handle("/", http.FileServer(http.Dir("dist/")))
		log.Fatal(http.ListenAndServe("localhost:"+Cfg.SRV.Port, nil))
	} else {
		os.Remove(Cfg.SRV.Unix)
		server := http.Server{}
		l, err := net.Listen("unix", Cfg.SRV.Unix)
		if err != nil {
			log.Println(err)
		}
		if err := os.Chmod(Cfg.SRV.Unix, 0777); err != nil {
			log.Fatal(err)
		}
		server.Serve(l)
	}

}
