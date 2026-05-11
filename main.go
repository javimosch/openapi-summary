package main
import ("encoding/json";"fmt";"os")
func main() {
	if len(os.Args) < 2 { fmt.Fprintln(os.Stderr,"Usage: openapi-summary <file.json>"); os.Exit(1) }
	for fi, path := range os.Args[1:] {
		data, err := os.ReadFile(path)
		if err != nil { fmt.Fprintln(os.Stderr,"Error:",err); continue }
		var spec map[string]any
		if err := json.Unmarshal(data, &spec); err != nil { fmt.Fprintln(os.Stderr,"JSON error:",err); continue }
		if fi > 0 { fmt.Println("---") }
		title, _ := spec["info"].(map[string]any)["title"].(string)
		version, _ := spec["info"].(map[string]any)["version"].(string)
		fmt.Printf(`{"file":"%s","title":"%s","version":"%s","endpoints":{`, path, title, version)
		endpoints, _ := spec["paths"].(map[string]any)
		methods := map[string]int{}
		type ep struct{ method, path, summary string }
		var eps []ep
		for p, item := range endpoints {
			im, _ := item.(map[string]any)
			for _, m := range []string{"get","post","put","patch","delete","options","head"} {
				if op, ok := im[m].(map[string]any); ok {
					methods[m]++
					s, _ := op["summary"].(string)
					eps = append(eps, ep{m, p, s})
				}
			}
		}
		first := true
		for _, m := range []string{"get","post","put","patch","delete","options","head"} {
			if methods[m] > 0 { if !first { fmt.Print(",") }; first = false; fmt.Printf(`"%s":%d`, m, methods[m]) }
		}
		fmt.Printf(`},"routes":[`)
		for i, e := range eps {
			if i > 0 { fmt.Print(",") }
			fmt.Printf(`{"%s":"%s","summary":"%s"}`, e.method, e.path, e.summary)
		}
		fmt.Println("]}")
	}
}
