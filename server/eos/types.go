package eos

type EosOptimize struct {
	Ems      EosEms     `json:"ems"`
	Battery  EosBattery `json:"pv_akku,omitempty"`
	Inverter Inverter   `json:"inverter,omitempty"`
	Vehicle  EosVehicle `json:"eauto,omitempty"`
	// temperature forecast in °C
	TempForecast []int `json:"temperature_forecast,omitempty"`
	// solution of the last optimization
	StartSolution []int `json:"start_solution,omitempty"`
}

type EosEms struct {
	// pv forecast in wh
	PVForecast []float64 `json:"pv_prognose_wh"`
	// price for grid consumption per wh
	GridTariff []float64 `json:"strompreis_euro_pro_wh"`
	// price for feed in per wh
	FeedInTariff []float64 `json:"einspeiseverguetung_euro_pro_wh"`
	// price for energy stored in the battery per wh
	BatteryTariff float64 `json:"preis_euro_pro_wh_akku"`
	// consumption of the house in wh without ev charging
	HouseConsumption []float64 `json:"gesamtlast"`
}

type EosBattery struct {
	// capacity in wh
	Capacity int `json:"capacity_wh"`
	// charging efficiency as float
	ChargingEfficiency float64 `json:"charging_efficiency,omitempty"`
	// discharging efficiency as float
	DischargingEfficiency float64 `json:"discharging_efficiency,omitempty"`
	// max charging power in w
	MaxPower int `json:"max_charge_power_w,omitempty,omitempty"`
	// soc in % at the start of the current hour // ToDO: right now current soc is used
	InitialSoc int `json:"initial_soc_percentage"`
	// min SOC in %
	MinSoc int `json:"min_soc_percentage,omitempty"`
	// max SOC in %
	MaxSoc int `json:"max_soc_percentage,omitempty"`
}

type EosVehicle struct {
	// capacity in wh
	Capacity int `json:"capacity_wh"`
	// Ladeeffizienz
	ChargingEfficiency float64 `json:"charging_efficiency,omitempty"`
	// max charging power in w
	MaxPower int `json:"max_charge_power_w,omitempty"`
	// current soc in %
	InitialSoc int `json:"initial_soc_percentage"`
	// min SOC in %
	MinSoc int `json:"min_soc_percentage,omitempty"`
	// max SOC in %
	MaxSoc int `json:"max_soc_percentage,omitempty"`
}

type Inverter struct {
	// max charging power of inverter in Wh
	MaxPower int `json:"max_power_wh"`
}

type EosOptimizeResponse struct {
	ACCharge              []float64 `json:"ac_charge"`
	DCCharge              []float64 `json:"dc_charge"`
	DischargeAllowed      []int     `json:"discharge_allowed"`
	EAutoChargeHoursFloat []float64 `json:"eautocharge_hours_float"`
	Result                struct {
		LastWhProStunde        []float64 `json:"Last_Wh_pro_Stunde"`
		EAutoSoCProStunde      []float64 `json:"EAuto_SoC_pro_Stunde"`
		EinnahmenEuroProStunde []float64 `json:"Einnahmen_Euro_pro_Stunde"`
		GesamtVerluste         float64   `json:"Gesamt_Verluste"`
		GesamtbilanzEuro       float64   `json:"Gesamtbilanz_Euro"`
		GesamteinnahmenEuro    float64   `json:"Gesamteinnahmen_Euro"`
		GesamtkostenEuro       float64   `json:"Gesamtkosten_Euro"`
		HomeApplianceWhPerHour []float64 `json:"Home_appliance_wh_per_hour"`
		KostenEuroProStunde    []float64 `json:"Kosten_Euro_pro_Stunde"`
		NetzbezugWhProStunde   []float64 `json:"Netzbezug_Wh_pro_Stunde"`
		NetzeinspeisungWh      []float64 `json:"Netzeinspeisung_Wh_pro_Stunde"`
		VerlusteProStunde      []float64 `json:"Verluste_Pro_Stunde"`
		AkkuSocProStunde       []float64 `json:"akku_soc_pro_stunde"`
		ElectricityPrice       []float64 `json:"Electricity_price"`
	} `json:"result"`
	EAutoObj struct {
		ChargeArray           []int   `json:"charge_array"`
		DischargeArray        []int   `json:"discharge_array"`
		DischargingEfficiency float64 `json:"discharging_efficiency"`
		Hours                 int     `json:"hours"`
		CapacityWh            int     `json:"capacity_wh"`
		ChargingEfficiency    float64 `json:"charging_efficiency"`
		MaxChargePowerW       int     `json:"max_charge_power_w"`
		SocWh                 int     `json:"soc_wh"`
		InitialSocPercentage  int     `json:"initial_soc_percentage"`
	} `json:"eauto_obj"`
	StartSolution []int `json:"start_solution"`
	WashingStart  int   `json:"washingstart"`
}

type last struct {
	Time string  `json:"time"`
	Last float64 `json:"Last"`
}

type eosGesamtlast struct {
	YearEnergy   int    `json:"year_energy"`
	MeasuredData []last `json:"measured_data"`
	Hours        int    `json:"hours"`
}
