package cis

var (
	iam1 = Recommendation{
		Name:   "Ensure that corporate login credentials are used instead of Gmail accounts",
		CisID:  "1.1",
		Scored: true,
		Level:  1,
	}
	iam2 = Recommendation{
		Name:   "Ensure that multi-factor authentication is enabled for all non- service accounts",
		CisID:  "1.2",
		Scored: false,
		Level:  1,
	}
	iam3 = Recommendation{
		Name:   "Ensure that there are only GCP-managed service account keys for each service account",
		CisID:  "1.3",
		Scored: true,
		Level:  1,
	}
	iam4 = Recommendation{
		Name:   "Ensure that ServiceAccount has no Admin privileges",
		CisID:  "1.4",
		Scored: true,
		Level:  1,
	}
	iam5 = Recommendation{
		Name:   "Ensure that IAM users are not assigned Service Account User role at project level",
		CisID:  "1.5",
		Scored: true,
		Level:  1,
	}
	iam6 = Recommendation{
		Name:   "Ensure user-managed/external keys for service accounts are rotated every 90 days or less",
		CisID:  "1.6",
		Scored: true,
		Level:  1,
	}
	iam7 = Recommendation{
		Name:   "Ensure that Separation of duties is enforced while assigning service account related roles to users",
		CisID:  "1.7",
		Scored: false,
		Level:  2,
	}
	iam8 = Recommendation{
		Name:   "Ensure Encryption keys are rotated within a period of 365 days",
		CisID:  "1.8",
		Scored: true,
		Level:  1,
	}
	iam9 = Recommendation{
		Name:   "Ensure that Separation of duties is enforced while assigning KMS related roles to users",
		CisID:  "1.9",
		Scored: true,
		Level:  2,
	}
	iam10 = Recommendation{
		Name:   "Ensure API keys are not created for a project",
		CisID:  "1.10",
		Scored: false,
		Level:  2,
	}
	iam11 = Recommendation{
		Name:   "Ensure API keys are restricted to use by only specified Hosts and Apps",
		CisID:  "1.11",
		Scored: false,
		Level:  1,
	}
	iam12 = Recommendation{
		Name:   "Ensure API keys are restricted to only APIs that application needs access",
		CisID:  "1.12",
		Scored: false,
		Level:  1,
	}
	iam13 = Recommendation{
		Name:   "Ensure API keys are rotated every 90 days",
		CisID:  "1.13",
		Scored: true,
		Level:  1,
	}
)
