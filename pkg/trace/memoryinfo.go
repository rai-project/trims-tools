package trace

import (
	"encoding/csv"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type MemoryInformation []*Trace

func (m MemoryInformation) Write(fmt string, output io.Writer) error {

	switch fmt {
	case "csv", "tsv":
		return m.writeCSV(fmt, output)
	case "table":
		return m.writeTable(output)
	}
	return errors.Errorf("the format %s is not a valid output format for gpu memory information", fmt)
}

func (m MemoryInformation) dsvHeader() []string {
	return []string{
		"model_name",
		"memory_usage",
		//"allocation_sizes",
	}
}

func (m MemoryInformation) dsvRows() [][]string {
	rows := [][]string{}
	for _, tr := range m {
		if tr.OtherDataRaw == nil {
			continue
		}
		modelName := tr.OtherDataRaw.ModelName
		events := tr.TraceEvents

		cudaMallocUsage := []uint64{}
		for _, event := range events {
			if event.Category == "memory" && event.Name == "cudaMalloc" {
				args, ok := event.Args.(map[string]interface{})
				if !ok {
					log.WithField("args", event.Args).Error("failed to cast event args to a map string")
					continue
				}
				size, err := cast.ToUint64E(args["size"])
				if err != nil {
					log.WithField("args", event.Args).WithField("size", args["size"]).Error("failed to cast cudaMalloc size to uint64")
					continue
				}
				cudaMallocUsage = append(cudaMallocUsage, size)
			}
		}

		allocs := []string{}
		cudaMallocTotal := float64(0)
		for _, u0 := range cudaMallocUsage {
			u := float64(u0) / float64(1024*1024)
			cudaMallocTotal += u
			allocs = append(allocs, cast.ToString(u))
		}

		row := []string{
			modelName,
			cast.ToString(cudaMallocTotal),
			//	strings.Join(allocs, ";"),
		}
		rows = append(rows, row)
	}
	return rows
}

func (m MemoryInformation) writeCSV(fmt string, output io.Writer) error {
	w := csv.NewWriter(output)
	if fmt == "tsv" {
		w.Comma = '\t'
	}
	w.Write(m.dsvHeader())
	w.WriteAll(m.dsvRows())
	w.Flush()
	return nil
}

func (m MemoryInformation) writeTable(output io.Writer) error {
	w := tablewriter.NewWriter(output)
	w.SetHeader(m.dsvHeader())
	w.AppendBulk(m.dsvRows())
	w.Render()
	return nil
}
