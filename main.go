package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	app := &cli.App{
		Name: "piedpiper",
		Commands: []*cli.Command{
			{
				Name:  "encode",
				Usage: "encode FILE OUTPUT",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "path of file to compress",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Usage:    "path of file to write compressed data to",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					input := ctx.String("input")
					inputFile, err := os.Open(input)
					if err != nil {
						return fmt.Errorf("failed to open input file: %w", err)
					}

					compressedBytes, err := encode(inputFile)
					if err != nil {
						return err
					}

					ouput := ctx.String("output")
					outputFile, err := os.OpenFile(ouput, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						return fmt.Errorf("failed to open output file: %w", err)
					}

					_, err = outputFile.Write(compressedBytes)
					if err != nil {
						return fmt.Errorf("failed to write compressed data to output file: %w", err)
					}

					return nil
				},
			},
			{
				Name:  "decode",
				Usage: "decode FILE OUTPUT",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "path of compressed file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Usage:    "path to file to decompress data to",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					input := ctx.String("input")
					inputFile, err := os.Open(input)
					if err != nil {
						return fmt.Errorf("failed to open compressed file: %w", err)
					}

					output := ctx.String("output")
					outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						return fmt.Errorf("failed to open output file: %w", err)
					}

					if err := decode(inputFile, outputFile); err != nil {
						return fmt.Errorf("failed to decompress data: %w", err)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	return nil
}
