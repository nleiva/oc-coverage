package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/buger/jsonparser"
	xr "github.com/nleiva/xrgrpc"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	user              = getenv("USER")
	password          = getenv("PASSWORD")
	address           = getenv("ADDRESS")
	paths             = "input/oc-paths.json"
	certificate       = "input/ems.pem"
	timeout           = 200
	id          int64 = 900
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Panicf("%s environment variable not set.", name)
	}
	return v
}

func check(err error, message string) {
	if err != nil {
		log.Fatalf("could not %s: %s", message, err)
	}
	log.Println(message)
}

// TODO: Add logs
func main() {
	// Specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername(user),
		xr.WithPassword(password),
		xr.WithHost(address),
		xr.WithCert(certificate),
		xr.WithTimeout(timeout),
	)
	check(err, "read target device parameters")

	// Read OpenConfig YANG paths.
	jsonpaths, err := ioutil.ReadFile(paths)
	check(err, "read OpenConfig YANG paths")

	// Setup a connection to the target.
	connection, ctx, err := xr.Connect(*router)
	check(err, "setup a gRPC client connection")
	defer connection.Close()

	// Get raw config from the target.
	raw, err := xr.ShowCmdTextOutput(ctx, connection, "show run", id)
	check(err, "get raw config from the target")

	// Write response to output file
	trimConfig(&raw)
	check(err, "trim the raw config")
	err = ioutil.WriteFile("output/initial.cfg", []byte(raw), 0644)
	check(err, "write raw config to a file")

	// Get OC config from target.
	output, err := xr.GetConfig(ctx, connection, string(jsonpaths), id+1)
	check(err, "get OpenConfig config from the target")

	// Write response to output file
	b := []byte(output)
	jconfig, _, _, err := jsonparser.Get(b, "data")
	check(err, "parse JSON OpenConfig config")
	err = ioutil.WriteFile("output/openconfig.json", jconfig, 0644)
	check(err, "write OpenConfig config to a file")

	// Take the router to a base config
	base, err := ioutil.ReadFile("input/base.cfg")
	check(err, "read the base config")
	err = xr.CommitReplace(ctx, connection, string(base), "", id+2)
	check(err, "apply the base config on the target")

	// Apply OpenConfig output to the target
	oc, err := ioutil.ReadFile("output/openconfig.json")
	check(err, "read the OpenConfig output")
	_, err = xr.MergeConfig(ctx, connection, string(oc), id+3)
	check(err, "apply OpenConfig config on the target")

	// Get new raw config from target.
	output, err = xr.ShowCmdTextOutput(ctx, connection, "show run", id+4)
	check(err, "get new raw config from the target")

	// Write new response to output file
	trimConfig(&output)
	check(err, "trim the new raw config")
	err = ioutil.WriteFile("output/openconfig.cfg", []byte(output), 0644)
	check(err, "write new raw config to file")

	// Take the router back to its initial config
	initial, err := ioutil.ReadFile("output/initial.cfg")
	check(err, "read the initial config")
	err = xr.CommitReplace(ctx, connection, string(initial), "", id+5)
	check(err, "set the initial config on the target")

	// Show the diff between initial and OpenConfig applied config
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(raw, output, true)
	err = ioutil.WriteFile("output/diff.cfg", []byte(dmp.DiffPrettyText(diffs)), 0644)
	check(err, "generate diff file")
}

func trimConfig(c *string) error {
	// Black magic to remove unnecesary header in config output
	if len(*c) < 107 {
		return fmt.Errorf("lenght of \"%v\" is less than 107", *c)
	}
	*c = (*c)[107:len(*c)]
	return nil
}
