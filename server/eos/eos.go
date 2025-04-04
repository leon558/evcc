package eos

import (
	"fmt"
	"net/http"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/api/globalconfig"
	"github.com/evcc-io/evcc/core"
	"github.com/evcc-io/evcc/server"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
)

type eos struct {
	*request.Helper
	config    *globalconfig.EosConfig
	log       *util.Logger
	influx    *server.Influx
	status    status
	optimized EosOptimizeResponse
}

type status struct {
	running        bool
	runningSince   time.Time
	nextRun        time.Time
	lastRunSeconds int
}

type eosConfig struct {
	ServerHost string `json:"server_eos_host"`
	ServerPort int    `json:"server_eos_port"`
	ConfigPath string `json:"config_file_path"`
}

func NewEosClient(conf *globalconfig.EosConfig, influx *server.Influx) (*eos, error) {
	log := util.NewLogger("eos")

	e := &eos{
		Helper: request.NewHelper(log),
		config: conf,
		log:    log,
		influx: influx,
	}

	// check endpoint /v1/config
	req, err := request.New(http.MethodGet, "http://"+conf.URL+"/v1/config", nil, request.JSONEncoding)

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

	return e, nil
}

func (e *eos) Run(site *core.Site) {
	if err := e.run(site); err != nil {
		e.log.ERROR.Println("eos:", e)
	}
	for range time.Tick(180 * time.Second) {
		if err := e.run(site); err != nil {
			e.log.ERROR.Println("eos:", e)
		}
	}
}

func (e *eos) run(site *core.Site) error {
	inverter := Inverter{
		MaxPower: e.config.InverterMaxPower,
	}
	ems := e.getEMS(site, e.config.BatteryTariff)
	battery := e.getBattery(site)
	car := e.getVehicle(site)
	solution := e.optimized.StartSolution

	eosOptimize := EosOptimize{
		Ems:           *ems,
		Battery:       *battery,
		Inverter:      inverter,
		Vehicle:       *car,
		StartSolution: solution,
	}

	e.status.running = true
	e.status.runningSince = time.Now()
	body := request.MarshalJSON(eosOptimize)
	req, err := request.New(http.MethodPost, "http://"+e.config.URL+"/optimize", body, request.JSONEncoding)
	e.Timeout = 120 * time.Second
	if err != nil {
		e.log.ERROR.Println("request error: ", err.Error())
		return nil
	}
	var result EosOptimizeResponse
	err = e.DoJSON(req, &result)
	if err != nil {
		e.log.ERROR.Println("request error: ", err.Error())
		return nil
	}
	e.status.running = false
	e.status.lastRunSeconds = int(time.Since(e.status.runningSince).Seconds())
	e.log.INFO.Println("result: ", result)
	e.optimized = result
	return nil
}

// Change to use the site.vehicle for everything / compare limit of vehicle and loadpoint
// and use the lowest value
func (e *eos) getVehicle(site *core.Site) *EosVehicle {
	loadpoints := site.Loadpoints()
	if len(loadpoints) == 0 {
		e.log.WARN.Println("no loadpoint found")
		return nil
	}
	loadpoint := loadpoints[0]

	vehicle, err := site.Vehicles().ByName(loadpoint.GetVehicleRef())
	if err != nil {
		e.log.WARN.Println("no vehicle found")
		return nil
	}
	instance := vehicle.Instance()

	capacity := instance.Capacity()
	soc, err := instance.Soc()
	if err != nil {
		e.log.WARN.Println("no vehicle soc found")
		return nil
	}
	minSoc := vehicle.GetMinSoc()
	maxSoc := loadpoint.GetLimitSoc()
	if soc := vehicle.GetLimitSoc(); soc > 0 && soc < maxSoc {
		maxSoc = vehicle.GetLimitSoc()
	}
	maxPower := int(loadpoint.GetMaxCurrent()) * loadpoint.GetPhases() * 230
	return &EosVehicle{
		Capacity:           int(capacity * 1000),
		ChargingEfficiency: 0.9,
		MaxPower:           maxPower,
		InitialSoc:         int(soc),
		MinSoc:             minSoc,
		MaxSoc:             maxSoc,
	}

}

func (e *eos) getBattery(site *core.Site) *EosBattery {
	capacity := site.GetBatteryCapacity()
	soc := site.GetBatterySoc()
	if capacity == 0 {
		e.log.WARN.Println("no battery found")
		return nil
	}
	if soc < 0 {
		e.log.WARN.Println("no battery soc found")
		return nil
	}
	return &EosBattery{
		Capacity:              int(capacity * 1000),
		InitialSoc:            int(soc),
		MaxPower:              e.config.InverterMaxPower,
		ChargingEfficiency:    0.9,
		DischargingEfficiency: 0.9,
		MinSoc:                e.config.BatMinSoc,
		MaxSoc:                e.config.BatMaxSoc,
	}

}

func (e *eos) getEMS(site *core.Site, batteryTariff float64) *EosEms {
	pv := e.getPvPrediction(site)
	grid := e.getGridTariff(site)
	feedin := e.getFeedInTariff(site)
	consumption := e.getConsumptionPrediction(e.config.Consumption, 48)

	return &EosEms{
		PVForecast:       pv,
		GridTariff:       grid,
		FeedInTariff:     feedin,
		BatteryTariff:    batteryTariff / 1000,
		HouseConsumption: consumption,
	}
}

func (e *eos) getGridTariff(site *core.Site) []float64 {
	var result []float64
	tariff := site.GetTariff(api.TariffUsageGrid)
	if tariff != nil {
		rates, err := tariff.Rates()
		if err == nil {
			if len(rates) > 48 {
				rates = rates[0:48]
			}
			for _, v := range rates {
				result = append(result, v.Price/1000)
			}
		}
	}
	if len(result) != 48 {
		req, err := request.New(http.MethodGet, "http://"+e.config.URL+"/strompreis", nil, request.JSONEncoding)
		if err != nil {
			e.log.ERROR.Println("request error: ", err.Error())
			return nil
		}
		var rates []float64
		err = e.DoJSON(req, &rates)
		if err != nil {
			e.log.ERROR.Println("request error: ", err.Error())
			return nil
		}
		// add the values from rates to the result array so that it has 48 values
		// needed when the tariff is not available for the next day
		for i := len(result); i < 48; i++ {
			result = append(result, rates[i])
		}
	}
	return result
}

func (e *eos) getFeedInTariff(site *core.Site) []float64 {
	tariff := site.GetTariff(api.TariffUsageFeedIn)
	if tariff == nil {
		e.log.ERROR.Println("no feed in tariff")
		return nil
	}
	rates, err := tariff.Rates()
	if err != nil {
		e.log.ERROR.Println("could not get feed in tariff", err)
		return nil
	}
	var result []float64
	if len(rates) > 48 {
		rates = rates[0:48]
	}
	for _, v := range rates {
		result = append(result, v.Price/1000)
	}
	return result
}

func (e *eos) getPvPrediction(site *core.Site) []float64 {
	pv := site.GetTariff(api.TariffUsageSolar)
	if pv == nil {
		e.log.ERROR.Println("no pv forecast")
	}
	rates, err := pv.Rates()
	if err != nil {
		e.log.ERROR.Println("could not get pv forecast", err)
	}
	if len(rates) > 48 {
		rates = rates[0:48]
	}
	var result []float64
	for _, v := range rates {
		result = append(result, v.Price)
	}
	return result
}

func (e *eos) queryConsumptionHistory() ([]last, error) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.Add(-168 * time.Hour)

	query := fmt.Sprintf(`from(bucket: "evcc") 
		|> range(start: time(v: "%s"), stop: time(v: "%s")) 
		|> filter(fn: (r) => r._measurement == "homePower" and r._field == "value") 
		|> aggregateWindow(every: 1h, fn: mean, createEmpty: true)`,
		start.Format(time.RFC3339), end.Format(time.RFC3339))

	result, err := e.influx.QueryDB(query)
	if err != nil {
		return nil, err
	}

	var values []last
	for result.Next() {
		v := result.Record().Value()
		t := result.Record().Time().Format(time.RFC3339)
		if v == nil {
			v = 0.0
		}
		values = append(values, last{
			Time: t,
			Last: v.(float64),
		})
	}
	return values, nil
}

func (e *eos) getConsumptionPrediction(consumption int, hours int) []float64 {
	last, err := e.queryConsumptionHistory()
	if err != nil || len(last) == 0 {
		e.log.ERROR.Println("keine verbrauchswerte", err)
	}
	payload := eosGesamtlast{
		YearEnergy:   consumption,
		MeasuredData: last,
		Hours:        hours,
	}
	req, err := request.New(http.MethodPost, "http://"+e.config.URL+"/gesamtlast", request.MarshalJSON(payload), request.JSONEncoding)
	if err != nil {
		e.log.ERROR.Println("request error: ", err.Error())
		return nil
	}
	var result []float64
	err = e.DoJSON(req, &result)
	if err != nil {
		e.log.ERROR.Println("request error: ", err.Error())
		return nil
	}
	return result
}
