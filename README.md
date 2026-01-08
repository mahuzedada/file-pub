# File Pub - VPC Testing Application

A Go application designed to test DevOps engineers' VPC and security group setup skills. This public file upload application requires proper configuration of VPC networking, RDS database connectivity, S3 access, and security groups.

## Overview

File Pub is a simple image upload and gallery application that tests the following AWS infrastructure components:

- **Public Subnet Configuration** - Web server accessible from the internet
- **Private Subnet Configuration** - RDS MySQL database in isolated subnet
- **Security Groups** - Proper ingress/egress rules
- **NAT Gateway** - EC2 instance accessing S3
- **IAM Roles** - Instance profile for S3 and RDS access
- **Network ACLs** - Optional additional security layer

## Architecture

```
Internet
    |
    v
Internet Gateway
    |
    v
Public Subnet (EC2 Instance)
    |
    +---> RDS MySQL (Private Subnet)
    |
    +---> S3 (via NAT Gateway or VPC Endpoint)
```

## Application Features

- Public web interface for image uploads
- Image gallery displaying all uploaded images
- Image metadata tracking (filename, size, type, upload time)
- Health check endpoint for connectivity testing
- Support for JPEG, PNG, GIF, and WebP images

## Prerequisites

### AWS Resources

1. **VPC** with public and private subnets
2. **RDS MySQL Instance** in private subnet
3. **S3 Bucket** for image storage
4. **EC2 Instance** in public subnet with IAM role
5. **Security Groups** properly configured

### Local Development

- Go 1.21 or higher
- Access to AWS account
- MySQL client (optional, for database setup)

## Setup Instructions

This application supports two environments:
- **Development**: Local machine with Docker MySQL and dev S3 bucket
- **Production**: AWS infrastructure with RDS and prod S3 bucket

### Development Setup (Local)

Perfect for testing the application locally before deploying to AWS.

#### 1. Clone Repository

```bash
git clone <repository-url>
cd file-pub
```

#### 2. Configure Development Environment

Edit `.env.dev` and set your existing S3 bucket name:

```bash
# Edit .env.dev
nano .env.dev
```

Required changes in `.env.dev`:
```bash
# Set your existing development S3 bucket name
S3_BUCKET=your-existing-dev-bucket

# Set your AWS credentials (if not using AWS CLI profile)
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

**Note**: S3 bucket should already exist. The setup script will not create it.

#### 3. Start Docker MySQL

This starts MySQL in detached mode:

```bash
make dev-setup
```

#### 4. Run Application in Development Mode

```bash
# Start MySQL and run application
make dev-run
```

Or use docker-compose for everything:

```bash
# Run both MySQL and app in Docker
make dev-compose
```

#### 5. Access Development Application

Open browser to:
- **Application**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`

#### Development Commands

```bash
make dev-setup         # Start Docker MySQL in detached mode
make dev-up           # Start Docker MySQL
make dev-down         # Stop Docker MySQL
make dev-run          # Run app with dev environment
make dev-logs         # View MySQL logs
make dev-compose      # Run everything with docker-compose
```

---

### Production Setup (AWS)

For deploying to AWS with RDS and production S3 bucket.

#### 1. Prerequisites

Ensure you have already created:
- RDS MySQL instance in private subnet
- S3 bucket for production
- VPC with proper security groups
- EC2 instance with IAM role

#### 2. Configure Production Environment

Edit `.env.prod` with your existing RDS and S3 bucket details:

```bash
# Edit .env.prod
nano .env.prod
```

Update `.env.prod`:
```bash
# Database Configuration (Your existing RDS endpoint)
DB_HOST=your-rds-endpoint.us-east-1.rds.amazonaws.com
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=your-secure-production-password
DB_NAME=filepub

# S3 Configuration (Your existing production bucket)
S3_BUCKET=your-existing-prod-bucket
S3_REGION=us-east-1
```

**Note**: RDS and S3 bucket should already exist. The setup script only validates configuration.

#### 3. Validate Configuration

```bash
make prod-setup
```

This validates your production configuration.

#### 4. Initialize Production Database

```bash
# Connect to RDS and run init script
mysql -h your-rds-endpoint -u admin -p < db/init.sql
```

#### 5. Build and Deploy to EC2

```bash
# Build production binary
make prod-build

# Deploy to EC2 (replace with your EC2 IP)
make prod-deploy SSH_HOST=ec2-user@1.2.3.4
```

#### 6. Start Application on EC2

SSH to your EC2 instance:

```bash
ssh ec2-user@your-ec2-ip

# Set environment variables
export $(cat .env | xargs)

# Run application
./file-pub
```

Or use systemd service (recommended):

```bash
# Create service file
sudo nano /etc/systemd/system/file-pub.service
```

```ini
[Unit]
Description=File Pub Application
After=network.target

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user
EnvironmentFile=/home/ec2-user/.env
ExecStart=/home/ec2-user/file-pub
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable file-pub
sudo systemctl start file-pub
sudo systemctl status file-pub
```

#### 7. Access Production Application

Open browser to:
- **Application**: `http://your-ec2-public-ip:8080`
- **Health Check**: `http://your-ec2-public-ip:8080/health`

#### Production Commands

```bash
make prod-setup        # Validate production configuration
make prod-build        # Build production binary
make prod-deploy       # Deploy to EC2
```

## VPC Setup Requirements

### Required VPC Configuration

#### 1. VPC Structure
- **VPC CIDR**: e.g., `10.0.0.0/16`
- **Public Subnet**: e.g., `10.0.1.0/24`
- **Private Subnet**: e.g., `10.0.2.0/24`
- **Internet Gateway**: Attached to VPC
- **NAT Gateway**: In public subnet (for S3 access)

#### 2. Route Tables

**Public Subnet Route Table:**
```
Destination      Target
10.0.0.0/16      local
0.0.0.0/0        igw-xxxxx (Internet Gateway)
```

**Private Subnet Route Table:**
```
Destination      Target
10.0.0.0/16      local
0.0.0.0/0        nat-xxxxx (NAT Gateway)
```

#### 3. Security Groups

**EC2 Security Group (Web Server):**
```
Inbound Rules:
- Type: HTTP, Port: 8080, Source: 0.0.0.0/0
- Type: SSH, Port: 22, Source: Your-IP/32

Outbound Rules:
- Type: All Traffic, Destination: 0.0.0.0/0
```

**RDS Security Group (Database):**
```
Inbound Rules:
- Type: MySQL/Aurora, Port: 3306, Source: EC2-Security-Group

Outbound Rules:
- Type: All Traffic, Destination: 0.0.0.0/0
```

#### 4. IAM Instance Profile

Create IAM role with policies:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::your-bucket-name"
    }
  ]
}
```

Attach this role to EC2 instance.

#### 5. S3 Bucket Policy

Configure your S3 bucket to allow public read access to uploaded images. This replaces the deprecated ACL approach.

In the AWS Console, go to your S3 bucket → Permissions → Bucket Policy and add:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::your-bucket-name/uploads/*"
    }
  ]
}
```

**Important Notes:**
- Replace `your-bucket-name` with your actual bucket name
- This policy allows public read access only to objects in the `uploads/` folder
- Ensure "Block all public access" is configured to allow this policy:
  - Uncheck "Block public access to buckets and objects granted through new public bucket or access point policies"
  - Or use the AWS CLI: `aws s3api put-public-access-block --bucket your-bucket-name --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=false,RestrictPublicBuckets=false"`

## Testing Checklist for Students

### Basic Connectivity Tests

- [ ] EC2 instance can be accessed from internet on port 8080
- [ ] Health check endpoint returns "OK" status
- [ ] Application can connect to RDS database
- [ ] Application can access S3 bucket
- [ ] Images can be uploaded successfully
- [ ] Uploaded images are displayed in gallery

### Security Tests

- [ ] RDS instance is NOT publicly accessible
- [ ] RDS security group only allows EC2 security group
- [ ] EC2 instance has no direct internet route to RDS
- [ ] S3 access works via NAT Gateway or VPC Endpoint
- [ ] SSH access is restricted to specific IP

### Advanced Tests

- [ ] VPC Flow Logs are enabled
- [ ] CloudWatch monitoring is configured
- [ ] Application logs are shipped to CloudWatch
- [ ] Auto Scaling Group is configured (optional)
- [ ] Load Balancer is configured (optional)
- [ ] HTTPS with SSL/TLS certificate (optional)

## Common Issues and Troubleshooting

### Issue: Cannot connect to database

**Check:**
1. RDS security group allows EC2 security group
2. Database endpoint is correct in `.env`
3. Database credentials are correct
4. RDS instance is in same VPC

**Debug:**
```bash
# Test from EC2 instance
mysql -h your-rds-endpoint -u admin -p

# Check security groups
aws ec2 describe-security-groups --group-ids sg-xxxxx
```

### Issue: Cannot upload to S3

**Check:**
1. IAM role is attached to EC2 instance
2. IAM role has S3 permissions
3. S3 bucket name is correct
4. EC2 can reach S3 (NAT Gateway or VPC Endpoint)

**Debug:**
```bash
# From EC2 instance
aws s3 ls s3://your-bucket-name

# Check IAM role
aws sts get-caller-identity
```

### Issue: Application not accessible from internet

**Check:**
1. Internet Gateway is attached to VPC
2. Public subnet route table points to IGW
3. EC2 security group allows port 8080
4. EC2 instance has public IP
5. Network ACLs allow traffic

**Debug:**
```bash
# Check routes
aws ec2 describe-route-tables --filters "Name=vpc-id,Values=vpc-xxxxx"

# Check security groups
aws ec2 describe-security-groups --filters "Name=vpc-id,Values=vpc-xxxxx"
```

### Issue: Health check fails

**Check:**
1. Database connectivity
2. S3 bucket accessibility
3. Application logs for errors

**Debug:**
```bash
# Check application logs
./file-pub 2>&1 | tee app.log

# Test health endpoint
curl http://localhost:8080/health
```

## API Endpoints

### GET /
- **Description**: Home page with upload form and image gallery
- **Response**: HTML page

### POST /upload
- **Description**: Upload image endpoint
- **Parameters**:
  - `image` (multipart/form-data): Image file
- **Accepted Types**: JPEG, PNG, GIF, WebP
- **Max Size**: 32 MB
- **Response**: Redirect to home page

### GET /health
- **Description**: Health check endpoint
- **Response**:
  - `200 OK`: All systems operational
  - `503 Service Unavailable`: Database or S3 issue

## Project Structure

```
file-pub/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Makefile                     # Build automation
├── Dockerfile                   # Container definition
├── docker-compose.yml           # Docker Compose for local dev
├── README.md                    # This file
├── .env.example                 # Environment template (general)
├── .env.dev                     # Development environment config
├── .env.prod                    # Production environment config
├── .gitignore                   # Git ignore rules
├── db/
│   └── init.sql                # Database schema
├── templates/
│   └── index.html              # HTML template
├── scripts/
│   ├── setup-dev.sh            # Development setup script
│   └── setup-prod.sh           # Production setup script
├── image/
│   ├── image_handler.go        # HTTP handlers
│   ├── image_service.go        # Business logic
│   ├── image_repository.go     # Database layer
│   ├── image_types.go          # Type definitions
│   └── image_errors.go         # Error definitions
└── internal/
    └── common/
        ├── validation.go       # Validation utilities
        ├── errors.go           # Error utilities
        ├── service.go          # Service utilities
        └── env.go              # Environment utilities
```

## Makefile Commands

### Development Commands (Local)
```bash
make dev-setup         # Start Docker MySQL in detached mode
make dev-up           # Start Docker MySQL
make dev-down         # Stop Docker MySQL
make dev-run          # Run app with dev environment
make dev-logs         # View MySQL logs
make dev-compose      # Run everything with docker-compose
```

### Production Commands (AWS)
```bash
make prod-setup        # Validate production configuration
make prod-build        # Build production binary
make prod-deploy       # Deploy to EC2
```

### General Commands
```bash
make help          # Show available commands
make deps          # Download dependencies
make build         # Build application binary
make run           # Run application locally
make test          # Run tests
make clean         # Clean build artifacts
make fmt           # Format code
make lint          # Run linter
make docker-build  # Build Docker image
make docker-run    # Run in Docker
make docker-stop   # Stop Docker container
make db-init       # Show database initialization script
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_HOST` | RDS endpoint | Yes | localhost |
| `DB_PORT` | MySQL port | No | 3306 |
| `DB_USER` | Database user | Yes | root |
| `DB_PASSWORD` | Database password | Yes | password |
| `DB_NAME` | Database name | Yes | filepub |
| `S3_BUCKET` | S3 bucket name | Yes | - |
| `S3_REGION` | AWS region | No | us-east-1 |
| `PORT` | Application port | No | 8080 |

## License

MIT License - Feel free to use for educational purposes.

## Contributing

This is a teaching tool. Contributions welcome for:
- Additional testing scenarios
- Documentation improvements
- Bug fixes
- Security enhancements

## Support

For issues or questions, please open a GitHub issue or contact the course instructor.
