package cis

var (
	sql1 = Recommendation{
		Name:   "Ensure that Cloud SQL database instance requires all incoming connections to use SSL",
		CisID:  "6.1",
		Scored: true,
		Level:  1,
	}
	sql2 = Recommendation{
		Name:   "Ensure that Cloud SQL database Instances are not open to the world",
		CisID:  "6.2",
		Scored: true,
		Level:  1,
	}
	sql3 = Recommendation{
		Name:   "Ensure that MySql database instance does not allow anyone to connect with administrative privileges",
		CisID:  "6.3",
		Scored: true,
		Level:  1,
	}
	sql4 = Recommendation{
		Name:   "Ensure that MySQL Database Instance does not allows root login from any Host",
		CisID:  "6.4",
		Scored: true,
		Level:  1,
	}
)
