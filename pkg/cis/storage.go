package cis

var (
	storage1 = Recommendation{
		Name:   "Ensure that Cloud Storage bucket is not anonymously or publicly accessible",
		CisID:  "5.1",
		Scored: true,
		Level:  1,
	}
	storage2 = Recommendation{
		Name:   "Ensure that there are no publicly accessible objects in storage buckets",
		CisID:  "5.2",
		Scored: false,
		Level:  1,
	}
	storage3 = Recommendation{
		Name:   "Ensure that logging is enabled for Cloud storage buckets",
		CisID:  "5.3",
		Scored: true,
		Level:  1,
	}
)
