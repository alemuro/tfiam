# TFIAM - Terraform IAM policies hardening

Run `tfiam` to generate an IAM policy only with the required permissions by your stack, and avoid providing more than needed permissions. 

## Installation

```
go install github.com/alemuro/tfiam/cmd/tfiam@latest
```

## Usage

Go to the folder where the Terraform code is stored, run `terraform init`, and then execute:

```
tfiam
```

This will produce the AWS IAM policy and will get back to you through stdout.


## Contributors welcome

Feel free to open a PR. Check the backlog or issues before.
