package set

import mapset "github.com/deckarep/golang-set/v2"

func test() {
	required := mapset.NewSet[string]()
	required.Add("cooking")
	required.Add("english")
	required.Add("math")
	required.Add("biology")
}
