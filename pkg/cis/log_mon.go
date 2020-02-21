package cis

var (
	logmon1 = Recommendation{
		Name:   "Ensure that Cloud Audit Logging is configured properly across all services and all users from a project",
		CisID:  "2.1",
		Scored: true,
		Level:  1,
	}
	logmon2 = Recommendation{
		Name:   "Ensure that sinks are configured for all Log entries",
		CisID:  "2.2",
		Scored: true,
		Level:  1,
	}
	logmon3 = Recommendation{
		Name:   "Ensure that object versioning is enabled on log-buckets",
		CisID:  "2.3",
		Scored: true,
		Level:  1,
	}
	logmon4 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for Project Ownership assignments/changes",
		CisID:  "2.4",
		Scored: true,
		Level:  1,
	}
	logmon5 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for Audit Configuration Changes",
		CisID:  "2.5",
		Scored: true,
		Level:  1,
	}
	logmon6 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for Custom Role changes",
		CisID:  "2.6",
		Scored: true,
		Level:  2,
	}
	logmon7 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for VPC Network Firewall rule changes",
		CisID:  "2.7",
		Scored: true,
		Level:  1,
	}
	logmon8 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for VPC network route changes",
		CisID:  "2.8",
		Scored: true,
		Level:  1,
	}
	logmon9 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for VPC network changes",
		CisID:  "2.9",
		Scored: true,
		Level:  1,
	}
	logmon10 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for Cloud Storage IAM permission changes",
		CisID:  "2.10",
		Scored: true,
		Level:  1,
	}
	logmon11 = Recommendation{
		Name:   "Ensure log metric filter and alerts exists for SQL instance configuration changes",
		CisID:  "2.11",
		Scored: true,
		Level:  1,
	}
)
