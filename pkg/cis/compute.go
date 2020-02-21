package cis

var (
	compute1 = Recommendation{
		Name:   "Ensure that instances are not configured to use the default service account with full access to all Cloud APIs",
		CisID:  "4.1",
		Scored: true,
		Level:  1,
	}
	compute2 = Recommendation{
		Name:   "Ensure 'Block Project-wide SSH keys' enabled for VM instances",
		CisID:  "4.2",
		Scored: false,
		Level:  1,
	}
	compute3 = Recommendation{
		Name:   "Ensure oslogin is enabled for a Project",
		CisID:  "4.3",
		Scored: true,
		Level:  1,
	}
	compute4 = Recommendation{
		Name:   "Ensure 'Enable connecting to serial ports' is not enabled for VM Instance",
		CisID:  "4.4",
		Scored: true,
		Level:  1,
	}
	compute5 = Recommendation{
		Name:   "Ensure that IP forwarding is not enabled on Instances",
		CisID:  "4.5",
		Scored: true,
		Level:  1,
	}
	compute6 = Recommendation{
		Name:   "Ensure VM disks for critical VMs are encrypted with Customer-Supplied Encryption Keys (CSEK)",
		CisID:  "4.6",
		Scored: true,
		Level:  2,
	}
)
