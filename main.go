package main

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"os"

	"github.com/mdouchement/hdr"
	_ "github.com/mdouchement/hdr/rgbe"
	"github.com/mdouchement/hdr/tmo"
)

type tonemapper func(hdr.Image) tmo.ToneMappingOperator

var entries = []struct {
	Input      string
	Output     string
	Name       string
	Tonemapper tonemapper
}{
	{Input: "forest_path.hdr", Output: "forest_path-linear.jpeg", Name: "Linear",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewLinear(m) }},
	{Input: "forest_path.hdr", Output: "forest_path-logarithmic.jpeg", Name: "Logarithmic",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewLogarithmic(m) }},
	{Input: "forest_path.hdr", Output: "forest_path-drago03.jpeg", Name: "Drago '03 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultDrago03(m) }},
	{Input: "forest_path.hdr", Output: "forest_path-reinhard05.jpeg", Name: "Reinhard '05 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultReinhard05(m) }},
	{Input: "forest_path.hdr", Output: "forest_path-custom_reinhard05.jpeg", Name: "Custom Reinhard '05 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultCustomReinhard05(m) }},
	{Input: "forest_path.hdr", Output: "forest_path-icam06.jpeg", Name: "iCAM06 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultICam06(m) }},
	//
	{Input: "memorial_o876.hdr", Output: "memorial_o876-linear.jpeg", Name: "Linear",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewLinear(m) }},
	{Input: "memorial_o876.hdr", Output: "memorial_o876-logarithmic.jpeg", Name: "Logarithmic",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewLogarithmic(m) }},
	{Input: "memorial_o876.hdr", Output: "memorial_o876-drago03.jpeg", Name: "Drago '03 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultDrago03(m) }},
	{Input: "memorial_o876.hdr", Output: "memorial_o876-reinhard05.jpeg", Name: "Reinhard '05 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultReinhard05(m) }},
	{Input: "memorial_o876.hdr", Output: "memorial_o876-custom_reinhard05.jpeg", Name: "Custom Reinhard '05 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultCustomReinhard05(m) }},
	{Input: "memorial_o876.hdr", Output: "memorial_o876-icam06.jpeg", Name: "iCAM06 (default)",
		Tonemapper: func(m hdr.Image) tmo.ToneMappingOperator { return tmo.NewDefaultICam06(m) }},
}

func main() {
	generateLDR()
	generateReadme()
}

func generateLDR() {
	for _, e := range entries {
		fi, err := os.Open(e.Input)
		check(err)
		defer fi.Close()

		m, _, err := image.Decode(fi)
		check(err)

		fmt.Printf("File: %s - Performing %s TMO\n", e.Input, e.Name)
		if hdrm, ok := m.(hdr.Image); ok {
			t := e.Tonemapper(hdrm)
			m = t.Perform()
		}

		fo, err := os.Create(e.Output)
		check(err)
		defer fo.Close()

		err = jpeg.Encode(fo, m, &jpeg.Options{Quality: 100})
		check(err)
	}
}

// README.md
// var endpoint = "" // dev
var endpoint = "https://github.com/mdouchement/hdr_examples/" // prod

var readme = `# HDR Gallery Examples

This gallery is generated using [mdouchement/hdr](https://github.com/mdouchement/hdr) Golang's package.

HDR Image credits:
- **forest_path.hdr** © Rafał Mantiuk
- **memorial_o876.hdr** © 1997 by P. Debevec, J. Malik

## Gallery

{{range .entries}}
- {{.Name}} ‒ {{.Input}}

![{{.Name}}]({{$.endpoint}}{{.Output}})
{{end}}

## License

**MIT**
`

func generateReadme() {
	f, err := os.Create("README.md")
	check(err)
	defer f.Close()

	t := template.Must(template.New("readme").Parse(readme))
	err = t.Execute(f, map[string]interface{}{
		"endpoint": endpoint,
		"entries":  entries,
	})
	check(err)

	check(f.Sync())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
