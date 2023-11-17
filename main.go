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
		Name:                 "piedpiper",
		Usage:                "a huffman coding compression tool",
		UsageText:            "piedpiper OPTIONS...",
		EnableBashCompletion: true,

		ArgsUsage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Usage:    "path of input file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "path of output file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "decode",
				Usage: "decode input to output",
			},
		},
		Action: func(ctx *cli.Context) error {
			input := ctx.String("input")
			output := ctx.String("output")

			inputFile, err := os.Open(input)
			if err != nil {
				return fmt.Errorf("failed to open input file: %w", err)
			}

			outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("failed to open output file: %w", err)
			}

			dec := ctx.Bool("decode")
			if dec {
				if err := decode(inputFile, outputFile); err != nil {
					return fmt.Errorf("failed to decompress data: %w", err)
				}

				return nil
			}

			compressedBytes, err := encode(inputFile)
			if err != nil {
				return err
			}

			_, err = outputFile.Write(compressedBytes)
			if err != nil {
				return fmt.Errorf("failed to write compressed data to output file: %w", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	return nil
}
