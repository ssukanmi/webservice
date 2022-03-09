# webservice
a web api application


### Prerequisites:
1) Go v1.17
2) MySQL 8.0.27


### How to deploy:
1) Clone the repository
2) Edit <code>.env</code>file with your db credentials
3) Run the command <code>go run server.go</code> in the terminal from the repository directory


# Packer


### Prerequisites:
1) packer 1.7.10


### How to deploy:
1) Export packer log <code>export PACKER_LOG=1</code>
2) Build image
```
packer build \
    -ver 'aws_access_key=' \
    -ver 'aws_secret_key=' \
    -ver 'aws_region=us-east-1' \
    -ver 'subnet_id=subnet-0f960342a73f60dc2' \
    -ver 'source_ami=ami-033b95fb8079dc481' \
    ami.pkr.hcl
```
#### OR
```
packer build --var-file=variables.pkrvars.hcl ami.pkr.hcl
```
