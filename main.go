package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/JaredReisinger/iffdump/iff"
	"github.com/JaredReisinger/iffdump/quetzal"
)

type chunkHeader struct {
	typeID [4]byte
	length uint32
}

func main() {
	logger := log.New()
	verboseFlag := flag.Bool("verbose", false, "verbose output")

	flag.Parse()

	if *verboseFlag {
		logger.SetLevel(log.DebugLevel)
		logger.Debug("using verbose output")
	}

	if len(flag.Args()) == 0 {
		logger.Fatal("no input file specified")
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		logger.WithError(err).Error("open file")
		return
	}
	defer f.Close()

	decoders := make(map[iff.TypeID]iff.ChunkDecoder, 0)
	quetzal.RegisterDecoders(decoders)
	decoder := iff.NewDecoder(decoders)

	c, err := decoder.Decode(f, logger)
	if err != nil {
		logger.WithError(err).Error("iff decode")
		return
	}

	logger.WithFields(log.Fields{
		"typeID": c.TypeID(),
		"length": c.Length(),
	}).Info("decoded iff file")

	log.Info(c.String())

	log.Info("DONE!")
}
