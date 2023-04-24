package database

type Measurements struct {
	ID                int    `influx:"id"`
	TypeMesure        string `influx:"typemesure"`
	DesignationMesure string `influx:"designationMesure"`
	ClasseMesure      string `influx:"classeMesure"`
	StatuMesure       string `influx:"statuMesure"`
	CategorieMesure   string `influx:"categorieMesure"`
	MinimumMesure     string `influx:"minimumMesure"`
	Target            string `influx:"target"`
}
