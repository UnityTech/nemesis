package cis

var (
	gke1 = Recommendation{
		Name:   "Ensure Stackdriver Logging is set to Enabled on Kubernetes Engine Clusters",
		CisID:  "7.1",
		Scored: true,
		Level:  1,
	}
	gke2 = Recommendation{
		Name:   "Ensure Stackdriver Monitoring is set to Enabled on Kubernetes Engine Clusters",
		CisID:  "7.2",
		Scored: false,
		Level:  1,
	}
	gke3 = Recommendation{
		Name:   "Ensure Legacy Authorization is set to Disabled on Kubernetes Engine Clusters",
		CisID:  "7.3",
		Scored: true,
		Level:  1,
	}
	gke4 = Recommendation{
		Name:   "Ensure Master authorized networks is set to Enabled on Kubernetes Engine Clusters",
		CisID:  "7.4",
		Scored: false,
		Level:  1,
	}
	gke5 = Recommendation{
		Name:   "Ensure Kubernetes Clusters are configured with Labels",
		CisID:  "7.5",
		Scored: false,
		Level:  1,
	}
	gke6 = Recommendation{
		Name:   "Ensure Kubernetes web UI / Dashboard is disabled",
		CisID:  "7.6",
		Scored: true,
		Level:  1,
	}
	gke7 = Recommendation{
		Name:   "Ensure Automatic node repair is enabled for Kubernetes Clusters",
		CisID:  "7.7",
		Scored: true,
		Level:  1,
	}
	gke8 = Recommendation{
		Name:   "Ensure Automatic node upgrades is enabled on Kubernetes Engine Clusters nodes",
		CisID:  "7.8",
		Scored: true,
		Level:  1,
	}
	gke9 = Recommendation{
		Name:   "Ensure Container-Optimized OS (COS) is used for Kubernetes Engine Clusters Node image",
		CisID:  "7.9",
		Scored: false,
		Level:  2,
	}
	gke10 = Recommendation{
		Name:   "Ensure Basic Authentication is disabled on Kubernetes Engine Clusters",
		CisID:  "7.10",
		Scored: true,
		Level:  1,
	}
	gke11 = Recommendation{
		Name:   "Ensure Network policy is enabled on Kubernetes Engine Clusters",
		CisID:  "7.11",
		Scored: true,
		Level:  1,
	}
	gke12 = Recommendation{
		Name:   "Ensure Kubernetes Cluster is created with Client Certificate enabled",
		CisID:  "7.12",
		Scored: true,
		Level:  1,
	}
	gke13 = Recommendation{
		Name:   "Ensure Kubernetes Cluster is created with Alias IP ranges enabled",
		CisID:  "7.13",
		Scored: true,
		Level:  1,
	}
	gke14 = Recommendation{
		Name:   "Ensure PodSecurityPolicy controller is enabled on the Kubernetes Engine Clusters",
		CisID:  "7.14",
		Scored: false,
		Level:  1,
	}
	gke15 = Recommendation{
		Name:   "Ensure Kubernetes Cluster is created with Private cluster enabled",
		CisID:  "7.15",
		Scored: true,
		Level:  1,
	}
	gke16 = Recommendation{
		Name:   "Ensure Private Google Access is set on Kubernetes Engine Cluster Subnets",
		CisID:  "7.16",
		Scored: true,
		Level:  1,
	}
	gke17 = Recommendation{
		Name:   "Ensure default Service account is not used for Project access in Kubernetes Clusters",
		CisID:  "7.17",
		Scored: true,
		Level:  1,
	}
	gke18 = Recommendation{
		Name:   "Ensure Kubernetes Clusters created with limited service account Access scopes for Project access",
		CisID:  "7.18",
		Scored: true,
		Level:  1,
	}
)
