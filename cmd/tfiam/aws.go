package main

type AWSStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}

type AWSPolicy struct {
	Version   string       `json:"Version"`
	Statement AWSStatement `json:"Statement"`
}

func generateAWSPolicy(permissions []string) *AWSPolicy {
	var awspol AWSPolicy
	var awsstat AWSStatement

	// Generate statement
	awsstat = AWSStatement{}
	awsstat.Effect = "Allow"
	awsstat.Resource = "*"
	awsstat.Action = []string{}

	for _, perm := range permissions {
		awsstat.Action = append(awsstat.Action, perm)
	}

	awsstat.Action = unique(awsstat.Action)

	// Generate policy
	awspol = AWSPolicy{}
	awspol.Version = "2012-10-17"
	awspol.Statement = awsstat
	return &awspol
}
