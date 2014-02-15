package geohash

func Encode(lat, lon float64, precision int) string {
	//TODO
	return ""
}

func Decode(gh string) *Box {
	//TODO
	return nil
}

type Interval struct {
	Min, Max float64
}

type Box struct {
	Lat, Lon Interval
}
