package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	. "github.com/pierrre/geohash"
)

var (
	precision int
	round     bool
)

func init() {
	flag.IntVar(&precision, "precision", 0, "Precision")
	flag.BoolVar(&round, "round", true, "Round")
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		flag.Usage()
	}

	results, err := processArgs(flag.Args())
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(results, " "))
}

func processArgs(args []string) ([]string, error) {
	var results []string
	for _, arg := range flag.Args() {
		result, err := processArg(arg)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func processArg(arg string) (string, error) {
	if strings.Contains(arg, ",") {
		return processArgLatLon(arg)
	}
	return processArgGeohash(arg)
}

func processArgLatLon(arg string) (string, error) {
	latLon := strings.Split(arg, ",")
	if len(latLon) != 2 {
		return "", fmt.Errorf("'%s'' is not a valid location (lat,lon)", arg)
	}

	lat, err := strconv.ParseFloat(latLon[0], 64)
	if err != nil {
		return "", err
	}

	lon, err := strconv.ParseFloat(latLon[1], 64)
	if err != nil {
		return "", err
	}

	var gh string
	if precision > 0 {
		gh = Encode(lat, lon, precision)
	} else {
		gh = EncodeAuto(lat, lon)
	}
	return gh, nil
}

func processArgGeohash(arg string) (string, error) {
	box, err := Decode(arg)
	if err != nil {
		return "", err
	}

	var p Point
	if round {
		p = box.Round()
	} else {
		p = box.Center()
	}

	return fmt.Sprintf("%v,%v", p.Lat, p.Lon), nil
}
