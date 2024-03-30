package jobs

var JobStatusList []string = []string{
	"New ⭐️",
	"Completed & Shipped ✅",
	"Invoiced 🧾",
}

func GetStatusIndex(search string) int {
	for i, s := range JobStatusList {
		if s == search {
			return i
		}
	}

	return 0
}
