package main

import (
	"flag"
	"fmt"
	"github.com/pierrre/geohash"
	"strconv"
	"strings"
)

func main() {
	var precision int
	flag.IntVar(&precision, "precision", 0, "Precision")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	results, err := processArgs(flag.Args(), precision)
	if err != nil {
		panic(err)
	}

	fmt.Print(strings.Join(results, " "))
}

func processArgs(args []string, precision int) ([]string, error) {
	var results []string
	for _, arg := range flag.Args() {
		result, err := processArg(arg, precision)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func processArg(arg string, precision int) (string, error) {
	if strings.Contains(arg, ",") {
		return processArgLatLon(arg, precision)
	} else {
		return processArgGeohash(arg)
	}
}

func processArgLatLon(arg string, precision int) (string, error) {
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
		gh = geohash.Encode(lat, lon, precision)
	} else {
		gh = geohash.EncodeAuto(lat, lon)
	}
	return gh, nil
}

func processArgGeohash(arg string) (string, error) {
	box, err := geohash.Decode(arg)
	if err != nil {
		return "", err
	}

	round := box.Round()

	return fmt.Sprintf("%v,%v", round.Lat, round.Lon), nil
}
