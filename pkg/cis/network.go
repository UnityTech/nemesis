package cis

var (
	network1 = Recommendation{
		Name:   "Ensure the default network does not exist in a project",
		CisID:  "3.1",
		Scored: true,
		Level:  1,
	}
	network2 = Recommendation{
		Name:   "Ensure legacy networks does not exists for a project",
		CisID:  "3.2",
		Scored: true,
		Level:  1,
	}
	network3 = Recommendation{
		Name:   "Ensure that DNSSEC is enabled for Cloud DNS",
		CisID:  "3.3",
		Scored: false,
		Level:  1,
	}
	network4 = Recommendation{
		Name:   "Ensure that RSASHA1 is not used for key-signing key in Cloud DNS DNSSEC",
		CisID:  "3.4",
		Scored: false,
		Level:  1,
	}
	network5 = Recommendation{
		Name:   "Ensure that RSASHA1 is not used for zone-signing key in Cloud DNS DNSSEC",
		CisID:  "3.5",
		Scored: false,
		Level:  1,
	}
	network6 = Recommendation{
		Name:   "Ensure that SSH access is restricted from the internet",
		CisID:  "3.6",
		Scored: true,
		Level:  2,
	}
	network7 = Recommendation{
		Name:   "Ensure that RDP access is restricted from the internet",
		CisID:  "3.7",
		Scored: true,
		Level:  2,
	}
	network8 = Recommendation{
		Name:   "Ensure Private Google Access is enabled for all subnetwork in VPC Network",
		CisID:  "3.8",
		Scored: true,
		Level:  2,
	}
	network9 = Recommendation{
		Name:   "Ensure VPC Flow logs is enabled for every subnet in VPC Network",
		CisID:  "3.9",
		Scored: true,
		Level:  1,
	}
)
