// Copied and modified from: https://github.com/bcicen/ctop
// MIT License - Copyright (c) 2017 VektorLab

package config

var (
	GlobalParams = defaultParams
	GlobalSwitches = defaultSwitches
)

func Init(stateFilter string, tagsFilter []string) {
	GlobalParams = defaultParams
	GlobalParams = append(GlobalParams, Param{
		Key: "stateFiler",
		Val: stateFilter, 
		Lavel: "Filter on task state",
	})

	tags := make(map[string]string)
	for _, v := range tagsFilter {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			return fmt.Errorf("tags must be of the form: KEY=VALUE")
		}
		tags[parts[0]] = parts[1]
	}

	GlobalParams = append(GlobalParams, Param{
		Key: "tagsFiler",
		Val: tagsFilter, 
		Lavel: "Filter on task state",
	})
}
