// This file contains routine for treating command line parameters.
package main

import (
	"flag"
)

// contains parameters determined by command line params.
type cliParams struct {
	sourceFilename string
	targetFilename string
}

func cli(pCliParams *cliParams) {

	flag.StringVar(&(pCliParams.sourceFilename), "from", "source.properties", "source .properties file")
	flag.StringVar(&(pCliParams.targetFilename), "to", "target.properties", "target .properties file")

	flag.Parse()

}
