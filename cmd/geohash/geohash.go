/*
Geohash command-line.

# Usage

Encode lat/lon to geohash:

	geohash 48.86,2.35

	u09tvqx

Decode geohash to lat/lon:

	geohash u09tvqx

	48.86,2.35

Custom precision:

	geohash -precision=12 48.86,2.35

	u09tvqxnnuph

Don't round:

	geohash -round=false u09tvqx

	48.85963439941406,2.3503875732421875

Multiple arguments:

	geohash 35.691015,139.766014 u09tvqx

	xn77h3qe0pmt 48.86,2.35

Stdin:

	echo "u09tvqx" | geohash

	48.86,2.35
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pierrre/geohash"
)

var (
	flagPrecision int
	flagRound     bool
)

func init() {
	flag.IntVar(&flagPrecision, "precision", 0, "Precision")
	flag.BoolVar(&flagRound, "round", true, "Round")
	flag.Parse()
}

func main() {
	if err := processSwitch(); err != nil {
		panic(err)
	}
}

func processSwitch() error {
	if flag.NArg() > 0 {
		return processArgs()
	}
	return processStdin()
}

func processArgs() error {
	args := flag.Args()
	results := make([]string, 0, len(args))
	for _, arg := range args {
		result, err := processValue(arg)
		if err != nil {
			return err
		}
		results = append(results, result)
	}
	_, _ = fmt.Fprintln(os.Stdout, strings.Join(results, " "))
	return nil
}

func processStdin() error {
	first := true
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		result, err := processValue(scanner.Text())
		if err != nil {
			return err
		}
		if first {
			first = false
		} else {
			_, _ = fmt.Fprint(os.Stdout, " ")
		}
		_, _ = fmt.Fprint(os.Stdout, result)
	}
	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	return nil
}

func processValue(v string) (string, error) {
	if strings.Contains(v, ",") {
		return processLatLon(v)
	}
	return processGeohash(v)
}

func processLatLon(latLon string) (string, error) {
	latLonSplit := strings.Split(latLon, ",")
	if len(latLonSplit) != 2 {
		return "", fmt.Errorf("'%s'' is not a valid location (lat,lon)", latLon)
	}

	lat, err := strconv.ParseFloat(latLonSplit[0], 64)
	if err != nil {
		return "", fmt.Errorf("latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(latLonSplit[1], 64)
	if err != nil {
		return "", fmt.Errorf("longitude: %w", err)
	}

	var gh string
	if flagPrecision > 0 {
		gh = geohash.Encode(lat, lon, flagPrecision)
	} else {
		gh = geohash.EncodeAuto(lat, lon)
	}
	return gh, nil
}

func processGeohash(arg string) (string, error) {
	box, err := geohash.Decode(arg)
	if err != nil {
		return "", fmt.Errorf("geohash: %w", err)
	}

	var p geohash.Point
	if flagRound {
		p = box.Round()
	} else {
		p = box.Center()
	}

	return fmt.Sprintf("%v,%v", p.Lat, p.Lon), nil
}
