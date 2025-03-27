package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
)

type eos struct {
	*request.Helper
	log       *util.Logger
	connected bool
}

type eosConfig struct {
	ServerHost string `json:"server_eos_host"`
	ServerPort int    `json:"server_eos_port"`
	ConfigPath string `json:"config_file_path"`
}

func NewEosClient(url string, consumption int) (*eos, error) {
	log := util.NewLogger("eos")

	e := &eos{
		Helper:    request.NewHelper(log),
		log:       util.NewLogger("eos"),
		connected: false,
	}

	// check endpoint /v1/config
	req, err := request.New(http.MethodGet, "http://"+url+"/v1/config", nil, request.JSONEncoding)

	if err != nil {
		// log.ERROR.Println(err)
		return nil, fmt.Errorf("failed requsting eos: %w", err)
	}

	var config eosConfig
	err = e.DoJSON(req, &config)

	if err != nil {
		// log.ERROR.Println(err)
		return nil, fmt.Errorf("failed getting config: %w", err)
	}

	e.log.INFO.Println("config:", config)
	return e, nil
}

func (e *eos) Run() {
	for range time.Tick(10 * time.Second) {
		if err := e.run(); err != nil {
			e.log.ERROR.Println("eos:", e)
		}
	}
}

func (e *eos) run() error {
	e.log.INFO.Println("eos:")
	return nil
}
