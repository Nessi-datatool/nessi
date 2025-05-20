# Cloud Integration Testing Plan

This document outlines the testing strategy for Nessi.dev's cloud integration features across AWS, Azure, and GCP.

## Test Environments

### AWS Test Environment
- **Account**: Development AWS account with S3 access
- **Regions**: Primary (us-west-2), Secondary (eu-central-1)
- **Authentication Methods**: IAM role, Access keys, AWS profile
- **Test Buckets**: 
  - `nessi-test-delta-tables` (contains sample Delta tables)
  - `nessi-test-results` (for storing test results)

### Azure Test Environment
- **Account**: Development Azure subscription
- **Regions**: Primary (West US 2), Secondary (North Europe)
- **Authentication Methods**: Storage account key, SAS token, Azure AD (using Azure SDK v1.6.1+)
- **Test Containers**:
  - `nessi-test-delta-tables` (contains sample Delta tables)
  - `nessi-test-results` (for storing test results)
- **SDK Version**: Azure SDK for Go v1.6.1 or newer

### GCP Test Environment
- **Account**: Development GCP project
- **Regions**: Primary (us-west1), Secondary (europe-west3)
- **Authentication Methods**: Service account key, Application default credentials
- **Test Buckets**:
  - `nessi-test-delta-tables` (contains sample Delta tables)
  - `nessi-test-results` (for storing test results)

## Test Data

### Sample Delta Tables
- **Small Table**: ~10MB, 1,000 rows, 10 columns
- **Medium Table**: ~100MB, 100,000 rows, 20 columns
- **Large Table**: ~1GB, 1,000,000 rows, 30 columns
- **Partitioned Table**: Partitioned by date, ~500MB total
- **Schema Evolution Table**: Multiple schema changes across versions

### Test Scenarios
1. **Basic Connectivity**: Verify basic connection to each cloud provider
2. **Authentication**: Test all authentication methods for each provider
3. **Delta Lake Operations**: Test all Delta Lake operations on cloud-stored tables
4. **Performance**: Measure performance with different table sizes
5. **Error Handling**: Test behavior with invalid credentials, missing tables, etc.
6. **Concurrency**: Test multiple simultaneous operations
7. **Network Issues**: Test behavior during network interruptions

## Test Cases

### 1. Connectivity Tests

| Test ID | Description | Expected Result |
|---------|-------------|-----------------|
| CONN-AWS-01 | Connect to AWS S3 with valid credentials | Successful connection |
| CONN-AWS-02 | Connect to AWS S3 with invalid credentials | Appropriate error message |
| CONN-AZURE-01 | Connect to Azure Blob Storage with valid credentials | Successful connection |
| CONN-AZURE-02 | Connect to Azure Blob Storage with invalid credentials | Appropriate error message |
| CONN-GCP-01 | Connect to GCP Cloud Storage with valid credentials | Successful connection |
| CONN-GCP-02 | Connect to GCP Cloud Storage with invalid credentials | Appropriate error message |

### 2. Delta Lake Operation Tests

| Test ID | Description | Expected Result |
|---------|-------------|-----------------|
| DELTA-AWS-01 | List Delta tables in AWS S3 | Correct list of tables |
| DELTA-AWS-02 | Get table schema from AWS S3 | Correct schema returned |
| DELTA-AWS-03 | Perform time travel on AWS S3 table | Correct version retrieved |
| DELTA-AZURE-01 | List Delta tables in Azure Blob Storage | Correct list of tables |
| DELTA-AZURE-02 | Get table schema from Azure Blob Storage | Correct schema returned |
| DELTA-AZURE-03 | Perform time travel on Azure Blob Storage table | Correct version retrieved |
| DELTA-GCP-01 | List Delta tables in GCP Cloud Storage | Correct list of tables |
| DELTA-GCP-02 | Get table schema from GCP Cloud Storage | Correct schema returned |
| DELTA-GCP-03 | Perform time travel on GCP Cloud Storage table | Correct version retrieved |

### 3. Performance Tests

| Test ID | Description | Expected Result |
|---------|-------------|-----------------|
| PERF-AWS-01 | Process small table in AWS S3 | Completes within 5 seconds |
| PERF-AWS-02 | Process medium table in AWS S3 | Completes within 30 seconds |
| PERF-AWS-03 | Process large table in AWS S3 | Completes within 3 minutes |
| PERF-AZURE-01 | Process small table in Azure Blob Storage | Completes within 5 seconds |
| PERF-AZURE-02 | Process medium table in Azure Blob Storage | Completes within 30 seconds |
| PERF-AZURE-03 | Process large table in Azure Blob Storage | Completes within 3 minutes |
| PERF-GCP-01 | Process small table in GCP Cloud Storage | Completes within 5 seconds |
| PERF-GCP-02 | Process medium table in GCP Cloud Storage | Completes within 30 seconds |
| PERF-GCP-03 | Process large table in GCP Cloud Storage | Completes within 3 minutes |

### 4. Error Handling Tests

| Test ID | Description | Expected Result |
|---------|-------------|-----------------|
| ERR-AWS-01 | Access non-existent table in AWS S3 | Appropriate error message |
| ERR-AWS-02 | Access table with insufficient permissions in AWS S3 | Permission error message |
| ERR-AZURE-01 | Access non-existent table in Azure Blob Storage | Appropriate error message |
| ERR-AZURE-02 | Access table with insufficient permissions in Azure Blob Storage | Permission error message |
| ERR-GCP-01 | Access non-existent table in GCP Cloud Storage | Appropriate error message |
| ERR-GCP-02 | Access table with insufficient permissions in GCP Cloud Storage | Permission error message |

## Automated Testing

### Unit Tests
- Create unit tests for each cloud provider implementation
- Mock cloud provider APIs for testing without actual cloud resources
- Test all error handling paths

### Integration Tests
- Create integration tests that connect to actual cloud resources
- Use test accounts with limited permissions
- Clean up test resources after tests complete

### Performance Tests
- Measure performance metrics for different operations
- Compare performance across cloud providers
- Identify bottlenecks and optimization opportunities

## Manual Testing Checklist

### AWS Manual Tests
- [ ] Configure AWS provider using CLI
- [ ] Configure AWS provider using configuration file
- [ ] Test IAM role authentication
- [ ] Test access key authentication
- [ ] List Delta tables in S3 bucket
- [ ] View schema of Delta table in S3
- [ ] Perform time travel on Delta table in S3
- [ ] Validate Delta table in S3 against rules

### Azure Manual Tests
- [ ] Configure Azure provider using CLI
- [ ] Configure Azure provider using configuration file
- [ ] Test account key authentication (using Azure SDK v1.6.1+)
- [ ] Test SAS token authentication (using Azure SDK v1.6.1+)
- [ ] Test Azure AD authentication (using Azure SDK v1.6.1+)
- [ ] List containers in Azure Blob Storage
- [ ] List blobs in Azure container
- [ ] Upload and download blobs to/from Azure Blob Storage
- [ ] List Delta tables in Azure Blob Storage
- [ ] View schema of Delta table in Azure Blob Storage
- [ ] Perform time travel on Delta table in Azure Blob Storage
- [ ] Validate Delta table in Azure Blob Storage against rules

### GCP Manual Tests
- [ ] Configure GCP provider using CLI
- [ ] Configure GCP provider using configuration file
- [ ] Test service account key authentication
- [ ] Test application default credentials
- [ ] List Delta tables in GCP Cloud Storage
- [ ] View schema of Delta table in GCP Cloud Storage
- [ ] Perform time travel on Delta table in GCP Cloud Storage
- [ ] Validate Delta table in GCP Cloud Storage against rules

## Test Execution Plan

1. **Phase 1**: Unit tests for all cloud providers
2. **Phase 2**: Integration tests with mock Delta tables
3. **Phase 3**: Performance tests with real Delta tables
4. **Phase 4**: Manual testing of CLI and user experience
5. **Phase 5**: Cross-cloud testing (e.g., comparing tables across providers)

## Test Reporting

- Generate test reports after each test run
- Track test coverage for cloud integration code
- Document any issues or limitations discovered during testing
- Create performance benchmarks for future comparison
