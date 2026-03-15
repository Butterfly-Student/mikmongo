# Mikhmon Examples

Real usage examples for Mikhmon library with actual MikroTik router connection.

## Router Configuration

**Default Router:**
- Host: 192.168.233.2:8728
- Username: admin
- Password: r00t

## Running Examples

```bash
cd examples/mikhmon
go run main.go
```

## Available Scenarios

1. **Voucher Generator** - Generate hotspot vouchers (vc/up mode)
2. **Profile Manager** - Create profiles with Mikhmon on-login script
3. **Multi Router** - Test multiple router connections
4. **Report Viewer** - View sales reports from /system/script

## Requirements

- MikroTik router accessible at 192.168.233.2:8728
- API service enabled on MikroTik
- User 'admin' with password 'r00t'
- Hotspot package installed on MikroTik

## Test Data Cleanup

All test data (vouchers, profiles) will be automatically cleaned up after testing.
